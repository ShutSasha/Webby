import React from 'react'

const SvgMock = React.forwardRef<SVGSVGElement, React.SVGProps<SVGSVGElement>>((props, ref) => (
  <svg ref={ref} data-testid="svg-mock" {...props} />
))

SvgMock.displayName = 'SvgMock'

export default SvgMock
