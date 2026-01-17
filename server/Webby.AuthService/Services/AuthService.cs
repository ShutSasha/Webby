using Microsoft.EntityFrameworkCore.ChangeTracking;
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

   public AuthService(IPasswordHasher passwordHasher, IRepository<User> repository, IMailService mailService)
   {
      _passwordHasher = passwordHasher;
      _repository = repository;
      _mailService = mailService;
   }

   public async Task Register(RegisterUserRequest request)
   {
      var candidate = await _repository.GetByPredicate(user => user.Email == request.Email);

      if (candidate == null)
      {
         throw new ApiException($"Email {request.Email} is already in use", 400);
      }

      var passwordHash = _passwordHasher.Generate(request.Password);
      var verificationCode = GenerateActivationCode();

      var user = new User
      {
         Username = request.Username,
         Email = request.Email,
         Password = passwordHash,
         About = string.Empty,
         //TODO: change to default user image from aws bucket
         AvatarUrl = "https://i.pinimg.com/originals/44/64/20/4464203a781eed3650f1fdd624c4d02a.jpg",
         VerificationCode = verificationCode,
         isVerified = false
      };

      await _mailService.SendVerificationCode(request.Email, verificationCode);
      await _repository.Add(user);
   }


   public async Task<User> Login(LoginUserRequest request)
   { 
      var user = (await _repository.GetByPredicate(user => user.Email == request.Email)).FirstOrDefault();

      if (user == null)
      {
         throw new ApiException("User not found", 404);
      }

      if (!user.isVerified)
      {
         throw new ApiException("User isn't verified", 400);
      }

      var isCorrectPassword = _passwordHasher.Verify(request.Password, user.Password);

      return isCorrectPassword switch
      {
         true => user,
         false => throw new ApiException("Incorrect password", 400)
      };
      
   }

   public async Task SendCode(ResendVerificationCodeRequest request)
   {
      var user = (await _repository.GetByPredicate(user => user.Email == request.Email)).FirstOrDefault();

      if (user == null)
      {
         throw new ApiException("User not found", 404);
      }

      var newVerificationCode = GenerateActivationCode();
      user.VerificationCode = newVerificationCode;

      await _repository.Update(user);
      await _mailService.SendVerificationCode(user.Email, newVerificationCode);
   }

   public async Task<User> VerifyEmail(VerifyUserRequest request)
   {
      var user = (await _repository.GetByPredicate(user => user.Email == request.Email)).FirstOrDefault();
      
      
      if (user == null || user.VerificationCode != request.VerificationCode)
      {
         throw new ApiException("Invalid email or verification code",400);
      }

      user.isVerified = true;
      user.VerificationCode = string.Empty;

      await _repository.Update(user);

      return user;
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