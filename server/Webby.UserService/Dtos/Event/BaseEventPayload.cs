using System.Text.Json.Serialization;

namespace Webby.UserService.Dtos.Event;

public record BaseEventPayload
{
   [JsonPropertyName("userId")]
   public Guid UserId { get; init; }

   [JsonPropertyName("value")]
   public int Value { get; init; } = 1;

   [JsonPropertyName("operation")] public EventOperation Operation { get; init; }
}

public enum EventOperation
{
   Increment = 1,
   Absolute
}