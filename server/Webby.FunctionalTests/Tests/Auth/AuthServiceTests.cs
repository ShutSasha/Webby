using System.Linq.Expressions;
using AutoMapper;
using FluentAssertions;
using NSubstitute;
using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Interfaces.Repositories;
using Webby.AuthService.Interfaces.Services;
using Webby.AuthService.Models;

namespace Webby.FunctionalTests.Tests.Auth;

public class AuthServiceTests
{
    private readonly IAuthRepository _repositoryMock;
    private readonly IPasswordHasher _passwordHasherMock;
    private readonly IMailService _mailServiceMock;
    private readonly IMapper _mapperMock;
    private readonly ITokenService _tokenServiceMock;
    private readonly IEventPublisher _eventPublisherMock;

    private readonly AuthService.Services.AuthService _sut;

    public AuthServiceTests()
    {
        _repositoryMock = Substitute.For<IAuthRepository>();
        _passwordHasherMock = Substitute.For<IPasswordHasher>();
        _mailServiceMock = Substitute.For<IMailService>();
        _mapperMock = Substitute.For<IMapper>();
        _tokenServiceMock = Substitute.For<ITokenService>();
        _eventPublisherMock = Substitute.For<IEventPublisher>();

        _sut = new AuthService.Services.AuthService(
            _passwordHasherMock,
            _repositoryMock,
            _mailServiceMock,
            _mapperMock,
            _tokenServiceMock,
            _eventPublisherMock
        );
    }

    #region Login Tests

    [Fact]
    public async Task Login_WhenUserNotFound_ShouldThrowApiException404()
    {
        // Arrange
        var request = new LoginUserRequest { Email = "test@test.com", Password = "password123" };
        
        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User>().AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.Login(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(404);
        exception.Which.Errors.Should().ContainKey("message");
    }

    [Fact]
    public async Task Login_WhenUserIsUnverified_ShouldThrowApiException400()
    {
        // Arrange
        var request = new LoginUserRequest { Email = "test@test.com", Password = "password123" };
        var user = new AuthService.Models.User { Email = "test@test.com", isVerified = false };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.Login(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(400);
        exception.Which.Errors.Should().ContainKey("isVerified");
    }

    [Fact]
    public async Task Login_WhenUserIsBanned_ShouldThrowApiException403()
    {
        // Arrange
        var request = new LoginUserRequest { Email = "test@test.com", Password = "password123" };
        var user = new AuthService.Models.User { Email = "test@test.com", isVerified = true, IsBanned = true };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.Login(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(403);
        exception.Which.Errors.Should().ContainKey("IsBanned");
    }

    [Fact]
    public async Task Login_WhenPasswordIsIncorrect_ShouldThrowApiException400()
    {
        // Arrange
        var request = new LoginUserRequest { Email = "test@test.com", Password = "wrong_password" };
        var user = new AuthService.Models.User { Email = "test@test.com", isVerified = true, IsBanned = false, Password = "hashed_password" };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());
            
        _passwordHasherMock.Verify(request.Password, user.Password).Returns(false);

        // Act
        Func<Task> act = async () => await _sut.Login(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(400);
        exception.Which.Errors.Should().ContainKey("password");
    }

    [Fact]
    public async Task Login_WhenValidCredentials_ShouldReturnLoginResponse()
    {
        // Arrange
        var request = new LoginUserRequest { Email = "test@test.com", Password = "correct_password" };
        var user = new AuthService.Models.User { Email = "test@test.com", isVerified = true, IsBanned = false, Password = "hashed_password" };
        var tokenExpiration = DateTime.UtcNow.AddHours(1);
        var expectedToken = new AuthToken { AccessToken = "jwt_token", ExpiresAt = new DateTimeOffset(tokenExpiration).ToUnixTimeSeconds() };
        var expectedUserDto = new UserDto { Email = "test@test.com" };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());
            
        _passwordHasherMock.Verify(request.Password, user.Password).Returns(true);
        _tokenServiceMock.GenerateToken(user).Returns(expectedToken);
        _mapperMock.Map<UserDto>(user).Returns(expectedUserDto);

        // Act
        var result = await _sut.Login(request);

        // Assert
        result.Should().NotBeNull();
        result.AccessToken.Should().Be(expectedToken.AccessToken);
        result.User.Email.Should().Be(expectedUserDto.Email);
    }

    #endregion

    #region Register Tests

    [Fact]
    public async Task Register_WhenUserAlreadyExistsAndVerified_ShouldThrowApiException400()
    {
        // Arrange
        var request = new RegisterUserRequest { Email = "exist@test.com", Username = "existUser", Password = "123" };
        var existingUser = new AuthService.Models.User { Email = "exist@test.com", Username = "existUser", isVerified = true };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { existingUser }.AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.Register(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(400);
        exception.Which.Errors.Should().ContainKey("email");
        exception.Which.Errors.Should().ContainKey("username");
    }

    [Fact]
    public async Task Register_WhenValidNewUser_ShouldAddUserAndSendEmail()
    {
        // Arrange
        var request = new RegisterUserRequest { Email = "new@test.com", Username = "newUser", Password = "password" };
        
        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User>().AsEnumerable());

        _passwordHasherMock.Generate(request.Password).Returns("hashed_password");

        // Act
        var result = await _sut.Register(request);

        // Assert
        result.Should().BeFalse();
        await _repositoryMock.Received(1).Add(Arg.Is<AuthService.Models.User>(u => 
            u.Email == request.Email && 
            u.Username == request.Username && 
            u.isVerified == false));
            
        await _mailServiceMock.Received(1).SendVerificationCode(request.Email, Arg.Any<string>());
    }

    #endregion

    #region PerformGoogleAuth Tests

    [Fact]
    public async Task PerformGoogleAuth_WhenUserIsBanned_ShouldThrowApiException403()
    {
        // Arrange
        var request = new GoogleAuthRequest { Email = "banned@google.com", Id = Guid.NewGuid() };
        var existingUser = new AuthService.Models.User { Email = "banned@google.com", IsBanned = true };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { existingUser }.AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.PerformGoogleAuth(request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Google auth error");
    }

    [Fact]
    public async Task PerformGoogleAuth_WhenNewUser_ShouldAddUserAndGenerateToken()
    {
        // Arrange
        var request = new GoogleAuthRequest { Email = "new@google.com", Id = Guid.NewGuid(), Name = "GoogleUser" };
        var tokenExpiration = DateTime.UtcNow.AddHours(1);
        var expectedToken = new AuthToken { AccessToken = "jwt_token", ExpiresAt = new DateTimeOffset(tokenExpiration).ToUnixTimeSeconds() };
        var expectedUserDto = new UserDto { Email = "new@google.com" };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User>().AsEnumerable());

        _tokenServiceMock.GenerateToken(Arg.Any<AuthService.Models.User>()).Returns(expectedToken);
        _mapperMock.Map<UserDto>(Arg.Any<AuthService.Models.User>()).Returns(expectedUserDto);

        // Act
        var result = await _sut.PerformGoogleAuth(request);

        // Assert
        result.Should().NotBeNull();
        result.AccessToken.Should().Be(expectedToken.AccessToken);
        
        await _repositoryMock.Received(1).Add(Arg.Is<AuthService.Models.User>(u => 
            u.Email == request.Email && 
            u.Username == request.Name && 
            u.isVerified == true)); 
    }

    #endregion

    #region VerifyEmail Tests

    [Fact]
    public async Task VerifyEmail_WhenCodeIsInvalid_ShouldThrowApiException400()
    {
        // Arrange
        var request = new VerifyUserRequest { Email = "test@test.com", VerificationCode = "000000" };
        var user = new AuthService.Models.User { Email = "test@test.com", VerificationCode = "123456" };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());

        // Act
        Func<Task> act = async () => await _sut.VerifyEmail(request);

        // Assert
        var exception = await act.Should().ThrowAsync<ApiException>();
        exception.Which.StatusCode.Should().Be(400);
        exception.Which.Errors.Should().ContainKey("code");
    }

    [Fact]
    public async Task VerifyEmail_WhenValid_ShouldSetVerifiedTrueAndClearCode()
    {
        // Arrange
        var request = new VerifyUserRequest { Email = "test@test.com", VerificationCode = "123456" };
        var user = new AuthService.Models.User { Email = "test@test.com", VerificationCode = "123456", isVerified = false };

        _repositoryMock.GetByPredicate(Arg.Any<Expression<Func<AuthService.Models.User, bool>>>())
            .Returns(new List<AuthService.Models.User> { user }.AsEnumerable());

        // Act
        await _sut.VerifyEmail(request);

        // Assert
        user.isVerified.Should().BeTrue();
        user.VerificationCode.Should().BeEmpty();
        await _repositoryMock.Received(1).Update(user);
    }

    #endregion
}