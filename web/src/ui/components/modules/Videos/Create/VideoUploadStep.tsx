import React from 'react'

type Props = {
  isUploading: boolean
  uploadProgress?: number
  onVideoSelect: (e: React.ChangeEvent<HTMLInputElement>) => void
}

export default function VideoUploadStep({ isUploading, uploadProgress = 0, onVideoSelect }: Props) {
  return (
    <div
      className="flex flex-col items-center justify-center py-20 border-2 border-dashed border-neutral-700 rounded-xl
        bg-neutral-800/20"
    >
      {isUploading ? (
        <div className="flex flex-col items-center gap-4 w-full max-w-sm px-4">
          <p className="text-foreground-tertiary font-semibold animate-pulse">Uploading video to server...</p>

          <div className="w-full flex flex-col gap-2 mt-2">
            <div className="flex justify-between items-center text-xs font-medium text-foreground-muted">
              <span>Progress</span>
              <span className="text-emerald-500">{uploadProgress}%</span>
            </div>
            <div className="h-2 w-full bg-neutral-800 rounded-full overflow-hidden">
              <div
                className="h-full bg-emerald-500 transition-all duration-300 ease-out"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
          </div>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-4">
          <div className="p-4 bg-neutral-800 rounded-full text-foreground-muted">
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
            <p className="text-foreground-tertiary font-semibold mb-1">Select video file to upload</p>
            <p className="text-foreground0 text-sm mb-6">MP4 or WebM</p>
          </div>

          <label
            className="cursor-pointer bg-emerald-500 hover:bg-emerald-600 text-foreground-inverse px-6 py-2.5 rounded-xl
              font-semibold transition-colors"
          >
            Choose File
            <input type="file" accept="video/mp4, video/webm" className="hidden" onChange={onVideoSelect} />
          </label>
        </div>
      )}
    </div>
  )
}
