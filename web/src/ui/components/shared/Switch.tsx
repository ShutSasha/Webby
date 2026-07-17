type Props = {
  isChecked: boolean
  toggle: () => void
}

export default function Switch({ isChecked, toggle }: Props) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={isChecked}
      onClick={toggle}
      className={` relative inline-flex p-1 w-16 cursor-pointer items-center rounded-full transition-colors duration-200
        ease-in-out focus:outline-emerald-500 outline-1 outline-transparent ring-0
        ${isChecked ? 'bg-emerald-500' : 'bg-background'} hover:bg-opacity-80 `}
    >
      <span
        className={` inline-block size-5 transform rounded-full shadow-lg transition duration-200 ease-in-out
          ${isChecked ? 'translate-x-9 bg-surface' : 'translate-x-0 bg-surface-tertiary'} `}
      />
    </button>
  )
}
