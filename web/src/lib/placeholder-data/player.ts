type TestVideo = {
  key: number
  type: 'mp4' | 'webm' | 'youtube'
  link: string
  description?: string
}

// sounds: '.ogg', '.ogv' -not supported
// not supprted: '.avi', '.flv', '.mkv'

const testVideoUrls: TestVideo[] = [
  {
    key: 1,
    type: 'mp4',
    link: 'https://www.shlomifish.org/Files/files/video/Stephen%20Colbert%20Sings%20Friday%20with%20Jimmy%20Fallon%20and%20The%20Roots%20(Late%20Night%20with%20Jimmy%20Fallon)-u_BszJzLyxI.mp4',
  },
  {
    key: 2,
    type: 'youtube',
    link: 'https://www.youtube.com/watch?v=Zmrj90wYt4c',
    description: '1 hour video',
  },
  {
    key: 3,
    type: 'webm',
    link: `https://www.shlomifish.org/Files/files/video/'Gravity'%20-%20Against%20The%20Current%20(Official%20Music%20Video)-S34KpOAgyTg.webm`,
  },
]

export { testVideoUrls }
