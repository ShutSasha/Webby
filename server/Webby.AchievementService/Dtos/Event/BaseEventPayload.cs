using System.Text.Json.Serialization;

namespace Webby.AchievementService.Dtos.Event;

public record BaseEventPayload
{
   [JsonPropertyName("userId")]
   public Guid UserId { get; init; }

   [JsonPropertyName("value")]
   public int Value { get; init; } = 1;

   [JsonPropertyName("is_increment_operation")] public bool IsIncrementOperation { get; init; }
}