using AutoMapper;
using FluentAssertions;
using Grpc.Core;
using Microsoft.Extensions.Logging;
using NSubstitute;
using NSubstitute.ExceptionExtensions;
using UserService.AchievementGrpcClient;
using Webby.NotificationService.GrpcClient;
using Webby.UserService.Clients;
using Webby.UserService.Dtos.Event;
using Webby.UserService.Dtos.Notification;
using Webby.UserService.Dtos.User;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Helpers;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models.Enums;

namespace Webby.FunctionalTests.Tests.User;

public class UserServiceTests
{
    private readonly IUserRepository _userRepositoryMock;
    private readonly IMapper _mapperMock;
    private readonly IStorageService _storageServiceMock;
    private readonly IUserPremiumRepository _userPremiumRepositoryMock;
    private readonly NotificationGrpcService.NotificationGrpcServiceClient _notificationGrpcClientMock;
    private readonly INotificationFactory _notificationFactoryMock;
    private readonly AchievementGrpcService.AchievementGrpcServiceClient _achievementGrpcClientMock;
    private readonly ILogger<UserService.Services.UserService> _loggerMock;
    private readonly VideoGrpcService.VideoGrpcServiceClient _videoGrpcClientMock;
    private readonly IComplaintRepository _complaintRepositoryMock;
    private readonly IEventPublisher _eventPublisherMock;

    private readonly UserService.Services.UserService _sut;

    public UserServiceTests()
    {
        _userRepositoryMock = Substitute.For<IUserRepository>();
        _mapperMock = Substitute.For<IMapper>();
        _storageServiceMock = Substitute.For<IStorageService>();
        _userPremiumRepositoryMock = Substitute.For<IUserPremiumRepository>();
        _notificationGrpcClientMock = Substitute.For<NotificationGrpcService.NotificationGrpcServiceClient>();
        _notificationFactoryMock = Substitute.For<INotificationFactory>();
        _achievementGrpcClientMock = Substitute.For<AchievementGrpcService.AchievementGrpcServiceClient>();
        _loggerMock = Substitute.For<ILogger<UserService.Services.UserService>>();
        _videoGrpcClientMock = Substitute.For<VideoGrpcService.VideoGrpcServiceClient>();
        _complaintRepositoryMock = Substitute.For<IComplaintRepository>();
        _eventPublisherMock = Substitute.For<IEventPublisher>();

        _sut = new UserService.Services.UserService(
            _userRepositoryMock,
            _mapperMock,
            _storageServiceMock,
            _userPremiumRepositoryMock,
            _notificationGrpcClientMock,
            _notificationFactoryMock,
            _achievementGrpcClientMock,
            _loggerMock,
            _videoGrpcClientMock,
            _complaintRepositoryMock,
            _eventPublisherMock
        );
    }

    #region GetUserInformation Tests

    [Fact]
    public async Task GetUserInformation_WhenUserNotFound_ShouldThrowApiException404()
    {
        // Arrange
        var userId = Guid.NewGuid();
        _userRepositoryMock.FindById(userId).Returns((UserService.Models.User)null);

        // Act
        Func<Task> act = async () => await _sut.GetUserInformation(userId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 404 && e.Message == "Get user information error");
    }

    [Fact]
    public async Task GetUserInformation_WhenGrpcFails_ShouldReturnEmptyAchievements()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var user = new UserService.Models.User { UserId = userId, Username = "TestUser" };
        var expectedUserDto = new UserDto { UserId = userId, Username = "TestUser" };

        _userRepositoryMock.FindById(userId).Returns(user);
        _mapperMock.Map<UserDto>(user).Returns(expectedUserDto);
        _userRepositoryMock.GetUserFollowBlock(userId).Returns(new UserService.Dtos.User.UserFollowStats());
        
        _achievementGrpcClientMock.GetPinnedAchievementsAsync(Arg.Any<GetPinnedAchievementsRequest>())
            .Throws(new RpcException(new Status(StatusCode.Unavailable, "Service offline")));

        // Act
        var result = await _sut.GetUserInformation(userId);

        // Assert
        result.Should().NotBeNull();
        result.User.Should().BeEquivalentTo(expectedUserDto);
        result.PinnedUserAchievements.Should().BeEmpty();
    }

    #endregion

    #region ProcessFollow Tests

    [Fact]
    public async Task ProcessFollow_WhenSelfFollow_ShouldThrowApiException400()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var request = new Webby.UserService.Dtos.User.UserFollowRequest { UserId = userId, FollowerId = userId };

        // Act
        Func<Task> act = async () => await _sut.ProcessFollow(request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 400 && e.Message == "Follow user error");
    }

    [Fact]
    public async Task ProcessFollow_WhenNotFollowing_ShouldAddFollowAndSendNotification()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var followerId = Guid.NewGuid();
        var request = new Webby.UserService.Dtos.User.UserFollowRequest { UserId = userId, FollowerId = followerId };
        var user = new Webby.UserService.Models.User { UserId = userId };
    
        var notificationDto = new SendNotificationDto 
        { 
            UserId = userId.ToString(), 
            Message = "New follower!",
            TargetIdentifier = followerId.ToString(),
            Title = "Follow Notification"
        };

        _userRepositoryMock.FindById(userId).Returns(user);
        _userRepositoryMock.HasUserFollow(userId, followerId).Returns(false); 
        _notificationFactoryMock.CreateNewFollowerNotification(userId, followerId).Returns(notificationDto);

        var grpcResponse = new CreateNotificationResponse(); 
        var asyncCall = new Grpc.Core.AsyncUnaryCall<CreateNotificationResponse>(
            Task.FromResult(grpcResponse), 
            Task.FromResult(new Grpc.Core.Metadata()), 
            () => Grpc.Core.Status.DefaultSuccess, 
            () => new Grpc.Core.Metadata(), 
            () => { });

        _notificationGrpcClientMock.SendNotificationToUserAsync(Arg.Any<CreateNotificationRequest>()).Returns(asyncCall);

        // Act
        var result = await _sut.ProcessFollow(request);

        // Assert
        result.Should().Be("Successfully followed user");
    
        // Оставляем await, так как это Task
        await _userRepositoryMock.Received(1).AddUserFollowing(userId, followerId);
    
        // УБИРАЕМ await, так как это AsyncUnaryCall от gRPC
        _notificationGrpcClientMock.Received(1).SendNotificationToUserAsync(Arg.Any<CreateNotificationRequest>());
    }

    #endregion

    #region BanUser Tests

    [Fact]
    public async Task BanUser_WhenRequesterIsModeratorAndTargetIsAdmin_ShouldThrowApiException403()
    {
        // Arrange
        var requestedUserId = Guid.NewGuid();
        var targetUserId = Guid.NewGuid();

        var requestedUser = new UserService.Models.User { UserId = requestedUserId, Role = Role.Moderator };
        var targetUser = new UserService.Models.User { UserId = targetUserId, Role = Role.Admin };

        _userRepositoryMock.FindById(requestedUserId).Returns(requestedUser);
        _userRepositoryMock.FindById(targetUserId).Returns(targetUser);

        // Act
        Func<Task> act = async () => await _sut.BanUser(requestedUserId, targetUserId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Ban user error");
    }

    [Fact]
    public async Task BanUser_WhenValid_ShouldUpdateUserAndCallGrpc()
    {
        // Arrange
        var requestedUserId = Guid.NewGuid();
        var targetUserId = Guid.NewGuid();

        var requestedUser = new UserService.Models.User { UserId = requestedUserId, Role = Role.Admin };
        var targetUser = new UserService.Models.User { UserId = targetUserId, Role = Role.User, IsBanned = false };

        _userRepositoryMock.FindById(requestedUserId).Returns(requestedUser);
        _userRepositoryMock.FindById(targetUserId).Returns(targetUser);
    
        var grpcResponse = new Google.Protobuf.WellKnownTypes.Empty();
        var asyncCall = new AsyncUnaryCall<Google.Protobuf.WellKnownTypes.Empty>(
            Task.FromResult(grpcResponse), Task.FromResult(new Metadata()), () => Status.DefaultSuccess, () => new Metadata(), () => { });
    
        _notificationGrpcClientMock.ReportBlockingAsync(Arg.Any<ReportBlockingRequest>()).Returns(asyncCall);

        // Act
        await _sut.BanUser(requestedUserId, targetUserId);

        // Assert
        targetUser.IsBanned.Should().BeTrue();
        
        await _userRepositoryMock.Received(1).Update(targetUser);
        await _complaintRepositoryMock.Received(1).SetIsBanComplaintStatus(targetUser.UserId, true);
        
        _notificationGrpcClientMock.Received(1).ReportBlockingAsync(
            Arg.Is<ReportBlockingRequest>(req => req != null && req.UserId == targetUserId.ToString()));
    }

    #endregion

    #region ChangeRole Tests

    [Fact]
    public async Task ChangeRole_WhenValidModeratorAssignment_ShouldPublishEvent()
    {
        // Arrange
        var requestedUserId = Guid.NewGuid();
        var targetUserId = Guid.NewGuid();
        var targetRole = Role.Moderator;

        var requestedUser = new UserService.Models.User { UserId = requestedUserId, Role = Role.Admin };
        var targetUser = new UserService.Models.User { UserId = targetUserId, Role = Role.User };

        _userRepositoryMock.FindById(requestedUserId).Returns(requestedUser);
        _userRepositoryMock.FindById(targetUserId).Returns(targetUser);

        var grpcResponse = new Google.Protobuf.WellKnownTypes.Empty();
        var asyncCall = new AsyncUnaryCall<Google.Protobuf.WellKnownTypes.Empty>(
            Task.FromResult(grpcResponse), Task.FromResult(new Metadata()), () => Status.DefaultSuccess, () => new Metadata(), () => { });
        
        _notificationGrpcClientMock.ReportBlockingAsync(Arg.Any<ReportBlockingRequest>()).Returns(asyncCall);

        // Act
        await _sut.ChangeRole(requestedUserId, targetUserId, targetRole);

        // Assert
        targetUser.Role.Should().Be(Role.Moderator);
        await _userRepositoryMock.Received(1).Update(targetUser);
        await _eventPublisherMock.Received(1).PublishAsync(Arg.Is<GetModeratorRoleEvent>(e => e.UserId == targetUserId));
    }

    #endregion
}