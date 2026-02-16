type TestVideo = {
  type: 'mp4' | 'webm' | 'avi' | 'flv' | 'mkv' | 'youtube'
  link: string
  description?: string
}

// sounds - '.ogg', '.ogv'

const testVideoUrls: TestVideo[] = [
  {
    type: 'mp4',
    link: 'https://www.shlomifish.org/Files/files/video/Stephen%20Colbert%20Sings%20Friday%20with%20Jimmy%20Fallon%20and%20The%20Roots%20(Late%20Night%20with%20Jimmy%20Fallon)-u_BszJzLyxI.mp4',
  },
  {
    type: 'youtube',
    link: 'https://www.youtube.com/watch?v=Zmrj90wYt4c',
    description: '1 hour video',
  },
]

export { testVideoUrls }
