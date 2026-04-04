import React from 'react'

type Props = {
  isUploading: boolean
  onVideoSelect: (e: React.ChangeEvent<HTMLInputElement>) => void
}

export default function VideoUploadStep({ isUploading, onVideoSelect }: Props) {
  return (
    <div
      className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-neutral-700 rounded-xl
        bg-neutral-800/20"
    >
      {isUploading ? (
        <div className="flex flex-col items-center gap-4">
          <div className="size-10 border-4 border-emerald-500/20 border-t-emerald-500 rounded-full animate-spin" />
          <p className="text-neutral-400 font-medium animate-pulse">Uploading video to server...</p>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-4">
          <div className="p-4 bg-neutral-800 rounded-full text-neutral-400">
            <svg className="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
          </div>
          <div className="text-center">
            <p className="text-neutral-200 font-semibold mb-1">Select video file to upload</p>
            <p className="text-neutral-500 text-sm mb-6">MP4, WebM or MKV</p>
          </div>

          <label
            className="cursor-pointer bg-emerald-500 hover:bg-emerald-600 text-neutral-950 px-6 py-2.5 rounded-xl
              font-semibold transition-colors"
          >
            Choose File
            <input type="file" accept="video/*" className="hidden" onChange={onVideoSelect} />
          </label>
        </div>
      )}
    </div>
  )
}
