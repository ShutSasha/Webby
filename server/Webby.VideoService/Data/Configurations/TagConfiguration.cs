using Amazon.S3.Model;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace Webby.VideoService.Data.Configurations;

public class TagConfiguration : IEntityTypeConfiguration<Models.Tag>
{
   public void Configure(EntityTypeBuilder<Models.Tag> builder)
   {
      builder.ToTable("Tags");

      builder.HasKey(t => t.TagId);

      builder.Property(t => t.Name)
         .IsRequired()
         .HasMaxLength(300);

      builder.HasIndex(t => t.Name)
         .IsUnique();

      builder.HasMany(t => t.VideoTags)
         .WithOne(vt => vt.Tag)
         .HasForeignKey(vt => vt.TagId);
   }
}
