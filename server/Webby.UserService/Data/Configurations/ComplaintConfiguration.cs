using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

namespace Webby.UserService.Data.Configurations;

using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using Webby.UserService.Models;

public class ComplaintConfiguration : IEntityTypeConfiguration<Complaint>
{
   public void Configure(EntityTypeBuilder<Complaint> builder)
   {
      builder.ToTable("Complaints");

      builder.HasKey(c => c.ComplaintId);

      builder.Property(c => c.ReasonType)
         .IsRequired()
         .HasMaxLength(200);

      builder.Property(c => c.AdditionalInfo)
         .HasMaxLength(2000);

      builder.Property(c => c.CreatedAt)
         .IsRequired();

      builder.Property(c => c.TargetType)
         .HasConversion<string>()
         .IsRequired();

      builder.HasIndex(c => new { c.TargetId, c.TargetType });
      
      builder.HasOne(c => c.Author)
         .WithMany(u => u.Complaints)
         .HasForeignKey(c => c.AuthorId)
         .OnDelete(DeleteBehavior.Cascade);
   }
}