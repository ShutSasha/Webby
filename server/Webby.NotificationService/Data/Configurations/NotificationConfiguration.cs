using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.NotificationService.Models;

namespace Webby.NotificationService.Data.Configurations;

public class NotificationConfiguration : IEntityTypeConfiguration<Notification>
{
   public void Configure(EntityTypeBuilder<Notification> builder)
   {
      builder.ToTable("Notifications");
      
      builder.Property(n => n.UserId)
         .IsRequired();

      builder.Property(n => n.NotificationStatus)
         .HasConversion<string>()
         .IsRequired();
      
      builder.HasKey(n => n.NotificationId);

      builder.HasIndex(n => n.UserId);
   }
}