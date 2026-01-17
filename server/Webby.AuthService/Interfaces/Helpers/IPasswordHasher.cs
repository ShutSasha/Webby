namespace Webby.AuthService.Interfaces.Helpers;

public interface IPasswordHasher
{
   string Generate(string password);
   bool Verify(string password, string hashPassword);
}