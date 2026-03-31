namespace Webby.VideoService.Dtos.External;

public class YouTubeVideoResponse
{
   public List<Item> Items { get; set; }
   public string? NextPageToken { get; set; }
   public PageInformation PageInfo { get; set; }

   public class Item
   {
      public dynamic Id { get; set; } 
      public Snippet Snippet { get; set; }
      public ContentDetails? ContentDetails { get; set; }
      public Statistics? Statistics { get; set; }
   }

   public class Snippet
   {
      public string Title { get; set; }
      public string Description { get; set; }
      public string ChannelId { get; set; }
      public string ChannelTitle { get; set; }
      public DateTime PublishedAt { get; set; }
      public Thumbnails Thumbnails { get; set; }
      public List<string>? Tags { get; set; }
   }

   public class ContentDetails { public string Duration { get; set; } }
   public class Statistics { public string ViewCount { get; set; } }
    
   public class Thumbnails { public Thumbnail Medium { get; set; } }
   public class Thumbnail { public string Url { get; set; } }

   public class PageInformation
   {
      public int TotalResults { get; set; }
      public int ResultsPerPage { get; set; }
   }
}