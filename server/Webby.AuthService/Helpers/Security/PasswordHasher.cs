using BCrypt.Net;
using Webby.AuthService.Interfaces.Helpers;

namespace Webby.AuthService.Helpers;

public class PasswordHasher : IPasswordHasher
{
   public string Generate(string password) =>
      BCrypt.Net.BCrypt.EnhancedHashPassword(password, HashType.SHA256);

   public bool Verify(string password, string? hashPassword) => 
      hashPassword != null && 
      BCrypt.Net.BCrypt.EnhancedVerify(password, hashPassword, HashType.SHA256);

}