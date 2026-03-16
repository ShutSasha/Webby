namespace Webby.VideoService.Helpers.Validation;

using System.ComponentModel.DataAnnotations;

public class MaxFileSizeAttribute : ValidationAttribute
{
   private readonly long _maxFileSize;

   public MaxFileSizeAttribute(long maxFileSize)
   {
      _maxFileSize = maxFileSize;
   }

   protected override ValidationResult? IsValid(object? value, ValidationContext validationContext)
   {
      if (value is IFormFile file)
      {
         if (file.Length > _maxFileSize)
         {
            return new ValidationResult(ErrorMessage ?? $"Maximum allowed file size is {_maxFileSize} bytes.");
         }
      }

      return ValidationResult.Success;
   }
}