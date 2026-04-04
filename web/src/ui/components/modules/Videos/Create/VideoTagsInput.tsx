import React, { MouseEvent, KeyboardEvent } from 'react'

type Props = {
  tags: string[]
  inputValue: string
  maxTags: number
  onInputChange: (val: string) => void
  onKeyDown: (e: KeyboardEvent<HTMLInputElement>) => void
  onRemoveTag: (e: MouseEvent<HTMLButtonElement>, tag: string) => void
}

export default function VideoTagsInput({ tags, inputValue, maxTags, onInputChange, onKeyDown, onRemoveTag }: Props) {
  return (
    <div className="flex flex-col gap-2 mt-2">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-neutral-300">Tags</span>
        <span className="text-xs text-neutral-500">
          {tags.length} / {maxTags}
        </span>
      </div>

      <div className="flex flex-col gap-3 p-3 bg-neutral-800 border border-neutral-700 rounded-lg transition-colors">
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {tags.map(tag => (
              <span
                key={tag}
                className="flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-neutral-700 text-xs font-medium
                  text-neutral-200"
              >
                #{tag}
                <button
                  type="button"
                  onClick={e => onRemoveTag(e, tag)}
                  className="text-neutral-400 hover:text-red-400 transition-colors focus:outline-none"
                  aria-label={`Remove tag ${tag}`}
                >
                  &times;
                </button>
              </span>
            ))}
          </div>
        )}

        <input
          type="text"
          value={inputValue}
          onChange={e => onInputChange(e.target.value)}
          onKeyDown={onKeyDown}
          disabled={tags.length >= maxTags}
          className="bg-transparent border-none text-sm text-neutral-100 focus:outline-none focus:ring-0 w-full
            placeholder:text-neutral-500 disabled:opacity-50 disabled:cursor-not-allowed"
          placeholder={tags.length >= maxTags ? 'Maximum tags reached' : 'Add a tag and press Enter or comma'}
        />
      </div>
    </div>
  )
}
