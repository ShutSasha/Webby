using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

public class PaymentConfiguration : IEntityTypeConfiguration<Payment>
{
   public void Configure(EntityTypeBuilder<Payment> builder)
   {
      builder.ToTable("Payments");

      builder.HasKey(p => p.PaymentId);
      
      builder.Property(p => p.Amount)
         .HasPrecision(18, 2)
         .IsRequired();

      builder.Property(p => p.Currency)
         .HasMaxLength(3)
         .IsRequired();
      
      builder.Property(p => p.ExternalId)
         .HasMaxLength(255);
      
      builder.HasIndex(p => p.ExternalId)
         .IsUnique();
      
      builder.HasIndex(p => p.UserId);

      builder.Property(p => p.Status)
         .HasConversion<string>()
         .IsRequired();

      builder.Property(p => p.CreatedAt)
         .HasDefaultValueSql("now()");
   }
}