using System.Text.Json;
using System.Text.Json.Serialization;

namespace Webby.VideoService.Helpers.Converters;

public class TimeSpanToStringConverter : JsonConverter<TimeSpan>
{
   public override TimeSpan Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
   {
      var value = reader.GetString();
      if (string.IsNullOrWhiteSpace(value)) 
         return TimeSpan.Zero;
      
      var parts = value.Split(':');
      if (parts.Length == 3 && 
          int.TryParse(parts[0], out int hours) &&
          int.TryParse(parts[1], out int minutes) &&
          int.TryParse(parts[2], out int seconds))
      {
         return new TimeSpan(hours, minutes, seconds);
      }
      
      return TimeSpan.TryParse(value, out var timeSpan) ? timeSpan : TimeSpan.Zero;
   }

   public override void Write(Utf8JsonWriter writer, TimeSpan value, JsonSerializerOptions options)
   {
      var formattedDuration = $"{(int)value.TotalHours:D2}:{value.Minutes:D2}:{value.Seconds:D2}";
      writer.WriteStringValue(formattedDuration);
   }
}