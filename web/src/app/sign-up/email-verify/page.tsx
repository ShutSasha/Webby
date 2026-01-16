'use client'

// eslint-disable-next-line import/named
import { OTPInput, SlotProps } from 'input-otp'

export default function VerifyPage() {
  const handleComplete = (code: string) => {
    console.log('code', code)
    // Server Action
  }

  return (
    <div className="flex flex-1 items-center justify-center">
      <div
        className="flex w-full bg-neutral-900 max-w-[700px] rounded-[20px] p-10 box-border flex-col items-center gap-6"
      >
        <div className="space-y-2 text-center">
          <h2 className="text-white text-2xl font-bold">Confirm email</h2>
          <p className="text-neutral-400 text-sm">We have sent a code to your email</p>
        </div>

        <OTPInput
          maxLength={6}
          onComplete={handleComplete}
          containerClassName="group flex items-center has-[:disabled]:opacity-50"
          render={({ slots }) => (
            <div className="flex gap-2">
              {slots.map((slot, idx) => (
                <Slot key={idx} {...slot} />
              ))}
            </div>
          )}
        />

        <button className="text-emerald-500 text-sm hover:underline mt-4">Resend code</button>
      </div>
    </div>
  )
}

function Slot(props: SlotProps) {
  return (
    <div
      className={` relative w-12 h-14 text-[20px] flex items-center justify-center transition-all duration-300 border-2
        rounded-xl
        ${props.isActive ? 'border-emerald-500 shadow-[0_0_10px_rgba(16,185,129,0.3)]' : 'border-neutral-700'}
        ${props.char ? 'text-white' : 'text-neutral-500'} `}
    >
      {props.char !== null ? <div>{props.char}</div> : null}

      {/* Caret (blinking bar) */}
      {props.isActive && props.char === null && (
        <div className="absolute inset-0 flex items-center justify-center animate-caret-blink">
          <div className="w-px h-8 bg-white" />
        </div>
      )}
    </div>
  )
}
