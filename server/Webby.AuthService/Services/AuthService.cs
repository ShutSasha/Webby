using AutoMapper;
using Microsoft.EntityFrameworkCore.ChangeTracking;
using MimeKit.Encodings;
using Webby.AuthService.Dtos;
using Webby.AuthService.Helpers.Exception;
using Webby.AuthService.Interfaces.Helpers;
using Webby.AuthService.Interfaces.Repositories;
using Webby.AuthService.Interfaces.Services;
using Webby.AuthService.Models;

namespace Webby.AuthService.Services;

public class AuthService : IAuthService
{
   private readonly IRepository<User> _repository;
   private readonly IPasswordHasher _passwordHasher;
   private readonly IMailService _mailService;
   private readonly IMapper _mapper;
   private readonly ITokenService _tokenService;

   public AuthService(IPasswordHasher passwordHasher, IRepository<User> repository,
      IMailService mailService, IMapper mapper, ITokenService tokenService)
   {
      _passwordHasher = passwordHasher;
      _repository = repository;
      _mailService = mailService;
      _mapper = mapper;
      _tokenService = tokenService;
   }

   public async Task<bool> Register(RegisterUserRequest request)
   {
      var existingUser = (await _repository
            .GetByPredicate(u => u.Email == request.Email || u.Username == request.Username))
         .FirstOrDefault();

      if (existingUser != null)
      {
         if (!existingUser.isVerified && 
             existingUser.Email == request.Email && 
             existingUser.Username == request.Username)
         {
            var newCode = GenerateActivationCode();

            existingUser.VerificationCode = newCode;
            existingUser.Password = _passwordHasher.Generate(request.Password);

            await _repository.Update(existingUser);
            await _mailService.SendVerificationCode(request.Email, newCode);

            return true;
         }

         var errors = new Dictionary<string, string>();

         if (existingUser.Email == request.Email)
            errors["email"] = "Email already in use";

         if (existingUser.Username == request.Username)
            errors["username"] = "Username already in use";

         throw new ApiException("Registration error", 400, errors);
      }

      var user = new User
      {
         Username = request.Username,
         Email = request.Email,
         Password = _passwordHasher.Generate(request.Password),
         About = string.Empty,
         //TODO: change to default user image from aws bucket
         AvatarUrl = "https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg",
         VerificationCode = GenerateActivationCode(),
         isVerified = false,
         Role = Role.User
      };

      await _repository.Add(user);
      await _mailService.SendVerificationCode(request.Email, user.VerificationCode);

      return false;
   }
   
   public async Task<LoginUserResponse> Login(LoginUserRequest request)
   {
      var loginUserResponse = new LoginUserResponse();
      
      var errors = new Dictionary<string, string>();
      var user = (await _repository.GetByPredicate(u => u.Email == request.Email)).FirstOrDefault();

      if (user == null)
      {
         errors["message"] = "User with specified credentials wasn't found";
         throw new ApiException("Login error", 404,errors);
      }
      
      if (!user.isVerified)
      {
         errors["isVerified"] = "User isn't verified";
         throw new ApiException("Login error", 400, errors);
      }

      if (!_passwordHasher.Verify(request.Password, user.Password))
      {
         errors["password"] = user.Password != null 
            ? "Incorrect password" 
            : "Missing password. Use a different login method and reset password for manually signing.";
         throw new ApiException("Login error", 400, errors);
      }

      var authTokenModel = await _tokenService.GenerateToken(user);
      
      loginUserResponse.User = _mapper.Map<UserDto>(user);
      loginUserResponse.AccessToken = authTokenModel.AccessToken;
      loginUserResponse.AccessTokenExpiresAt = authTokenModel.ExpiresAt;

      return loginUserResponse;
   }
   
   public async Task SendCode(ResendVerificationCodeRequest request)
   {
      var user = (await _repository.GetByPredicate(user => user.Email == request.Email)).FirstOrDefault();
      var errors = new Dictionary<string, string>();
      
      if (user == null)
      {
         errors["message"] = "User with specified email wasn't found";
         throw new ApiException("Send code error", 404,errors);
      }

      var newVerificationCode = GenerateActivationCode();
      user.VerificationCode = newVerificationCode;

      await _repository.Update(user);
      await _mailService.SendVerificationCode(user.Email, newVerificationCode);
   }

   public async Task VerifyEmail(VerifyUserRequest request)
   {
      var user = (await _repository.GetByPredicate(user => user.Email == request.Email)).FirstOrDefault();
      var errors = new Dictionary<string, string>();
      
      if (user == null)
      {
         errors["email"] = "Invalid email";
         
         throw new ApiException("Verification error",400, errors);
      }

      if (user.isVerified)
      {
         errors["code"] = "User is already verified";
         
         throw new ApiException("Verification error",400, errors);
      }

      if (user.VerificationCode != request.VerificationCode)
      {
         errors["code"] = user.VerificationCode == string.Empty ? "User is already verified" : "Invalid verification code";
         
         throw new ApiException("Verification error",400, errors);
      }

      user.isVerified = true;
      user.VerificationCode = string.Empty;

      await _repository.Update(user);
      
   }
   
   public async Task<LoginUserResponse> PerformGoogleAuth(GoogleAuthRequest request)
   {
      var user = (await _repository
            .GetByPredicate(u => u.UserId == request.Id || u.Email == request.Email))
         .FirstOrDefault();
      
      var loginUserResponse = new LoginUserResponse();
      AuthToken authToken;
      
      if (user != null)
      {
         authToken = await _tokenService.GenerateToken(user);
         
         loginUserResponse.User = _mapper.Map<UserDto>(user);
         loginUserResponse.AccessToken = authToken.AccessToken;
         loginUserResponse.AccessTokenExpiresAt = authToken.ExpiresAt;

         return loginUserResponse;
      }
      var newUser = new User
      {
         UserId = request.Id,
         Email = request.Email,
         Username = request.Name,
         AvatarUrl = request.Image,
         About = string.Empty,
         isVerified = true,
         VerificationCode = string.Empty,
         Password = null,
         Role = Role.User
      };

      await _repository.Add(newUser);
         
      authToken = await _tokenService.GenerateToken(newUser);
         
      loginUserResponse.User = _mapper.Map<UserDto>(newUser);
      loginUserResponse.AccessToken = authToken.AccessToken;
      loginUserResponse.AccessTokenExpiresAt = authToken.ExpiresAt;

      return loginUserResponse;
   }

   public async Task<LoginUserResponse> RefreshToken(string accessToken)
   {
      var errors = new Dictionary<string, string>();
      var loginUserResponse = new LoginUserResponse();

      var userId = await _tokenService.ExtractUserInfo(accessToken);

      var existingUser = await _repository.FindById(userId);

      if (existingUser == null)
      {
         errors["user"] = "User with specified id wasn't found";
         throw new ApiException("Refresh token error", 404, errors);
      }

      var authTokenModel = await _tokenService.GenerateToken(existingUser);
      
      loginUserResponse.User = _mapper.Map<UserDto>(existingUser);
      loginUserResponse.AccessToken = authTokenModel.AccessToken;
      loginUserResponse.AccessTokenExpiresAt = authTokenModel.ExpiresAt;

      return loginUserResponse;
   }

   public async Task<UserDto> ChangeUserPassword(Guid userId, ChangeUserPasswordRequest request)
   {
      var user = await _repository.FindById(userId);
      
      if (user == null)
      {
         throw new ApiException("Change password error", 400, "User with specified id wasn't found");
      }

      if (user.VerificationCode != string.Empty)
      {
         throw new ApiException("Change password error", 401, "User isn't verified");
      }

      if (request.CurrentPassword != null)
      {
         if (!_passwordHasher.Verify(request.CurrentPassword,user.Password))
         {
            throw new ApiException("Change password error", 400, "Current password is incorrect");
         }
      }
      
      var newPasswordHash = _passwordHasher.Generate(request.NewPassword);
      user.Password = newPasswordHash;

      await _repository.Update(user);

      return _mapper.Map<UserDto>(user);
   }

   private string GenerateActivationCode()
   {
      const int length = 6;
      var random = new Random();
      var code = new char[length];

      for (int i = 0; i < length; i++)
      {
         code[i] = (char)('0' + random.Next(0, 10));
      }

      return new string(code);
   }
}