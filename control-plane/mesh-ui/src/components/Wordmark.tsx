interface WordmarkProps {
  height?: number
  className?: string
}

// Pixel "MESH" brand wordmark. Recolorable via currentColor / text-* classes.
const RECTS: Array<[number, number, number, number]> = [
  [54, 171, 41, 164], [145, 208, 42, 127], [237, 208, 42, 127], [54, 171, 183, 37],
  [305, 168, 42, 166], [305, 168, 141, 37], [305, 232, 141, 37], [305, 297, 141, 37],
  [473, 168, 113, 37], [473, 168, 42, 101], [473, 232, 150, 37], [581, 269, 42, 65],
  [473, 297, 150, 37], [648, 168, 42, 166], [794, 168, 42, 166], [648, 232, 188, 37],
]

export function Wordmark({ height = 20, className = "" }: WordmarkProps) {
  return (
    <svg
      height={height}
      width={height * (782 / 167)}
      viewBox="54 168 782 167"
      fill="currentColor"
      shapeRendering="crispEdges"
      className={className}
      style={{ display: "block" }}
    >
      {RECTS.map((r, i) => (
        <rect key={i} x={r[0]} y={r[1]} width={r[2]} height={r[3]} />
      ))}
    </svg>
  )
}
