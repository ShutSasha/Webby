import SignUpForm from '@/ui/components/modules/SignUp/SignUpForm'

export default function Page() {
  return (
    <div className="flex flex-1 items-center justify-center">
      <div className="flex w-full bg-neutral-900 max-w-[460px] rounded-[20px] p-5 gap-4 box-border">
        <SignUpForm />
      </div>
    </div>
  )
}
