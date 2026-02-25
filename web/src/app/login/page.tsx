import LoginForm from '@/ui/components/modules/Login/LoginForm'

export default function Page() {
  return (
    <div className="flex flex-1 items-center justify-center">
      <div className="flex w-full bg-neutral-900 max-w-[460px] rounded-[20px] px-7 py-5 gap-4 box-border">
        <LoginForm />
      </div>
    </div>
  )
}
