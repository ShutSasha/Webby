declare module '*.svg' {
  import React from 'react'

  const content: React.FC<React.SVGProps<SVGSVGElement>>
  export default content

  export const ReactComponent: React.FC<React.SVGProps<SVGSVGElement>>
}

declare module '*.svg?url' {
  const content: string
  export default content
}
