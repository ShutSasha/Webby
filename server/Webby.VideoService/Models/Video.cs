using System.Diagnostics.CodeAnalysis;
using Webby.VideoService.Models.Enums;

namespace Webby.VideoService.Models;

public class Video
{
   public Video()
   {
   }

   [SetsRequiredMembers]
   public Video(bool isPrivate, Guid userId, string name)
   {
      IsPrivate = isPrivate;
      UserId = userId;
      Name = name;
   }

   public Guid VideoId { get; set; }
   public Guid UserId { get; set; }
   public required string Name { get; set; }
   public VideoStatus VideoUploadStatus { get; set; }
   public long Duration { get; set; }
   public string? Description { get; set; }
   public int Views { get; set; }
   public DateTime CreatedAt { get; set; }
   public string? VideoUrl { get; set; }
   public string? PreviewUrl { get; set; }
   public bool IsPrivate { get; set; }
   public bool IsPublished { get; set; }
   public ICollection<VideoTag>? VideoTags { get; set; }
   
}