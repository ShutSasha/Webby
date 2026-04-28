'use client'

import { PREMIUM_FEATURES } from '@/lib/placeholder-data/premium'
import { cn } from '@/lib/utils/general.utils'
import Button from '@/ui/components/shared/Button'

export default function PremiumPage() {
  return (
    <div className="relative min-h-screen w-full overflow-hidden pb-32 rounded-2xl">
      <div
        className="absolute top-[-10%] left-[-10%] w-[500px] h-[500px] rounded-full bg-emerald-500/10 blur-[120px]
          pointer-events-none"
      />
      <div
        className="absolute bottom-[20%] right-[-10%] w-[600px] h-[600px] rounded-full bg-emerald-500/5 blur-[150px]
          pointer-events-none"
      />

      <div className="max-w-6xl mx-auto px-6 md:px-12 pt-20">
        <section className="flex flex-col items-center text-center mb-24 md:mb-32">
          <h1 className="text-4xl md:text-5xl font-bold text-neutral-100 mb-6 tracking-tight">
            Webby Premium Subscription —<br /> your key to limitless viewing
          </h1>
          <p className="text-neutral-400 text-lg max-w-2xl mb-10 leading-relaxed">
            More comfort, freedom, and features for true co-watching enthusiasts.
          </p>
          <Button
            viewType="confirm"
            className="rounded-xl px-10 py-3.5 text-lg font-bold bg-emerald-500 hover:bg-emerald-400 text-neutral-900
              shadow-lg shadow-emerald-500/20 transition-all hover:scale-105"
          >
            Get Premium
          </Button>
        </section>

        <section className="flex flex-col gap-24 md:gap-32">
          {PREMIUM_FEATURES.map((feature, index) => {
            const isReversed = index % 2 !== 0

            return (
              <div
                key={feature.id}
                className={cn(
                  'flex flex-col md:flex-row items-center gap-12 md:gap-20',
                  isReversed ? 'md:flex-row-reverse' : '',
                )}
              >
                <div className="flex-1 flex flex-col items-start w-full">
                  <div className="flex items-center gap-3 mb-4">
                    <feature.icon className="w-8 h-8 text-emerald-500" />
                    <h2 className="text-2xl md:text-3xl font-bold text-neutral-100">{feature.title}</h2>
                  </div>

                  <p className="text-neutral-400 text-[15px] md:text-base leading-relaxed mb-6">
                    {feature.description}
                  </p>

                  <div
                    className="flex flex-col gap-3 mb-8 w-full p-5 rounded-2xl bg-neutral-900/50 border
                      border-neutral-800"
                  >
                    <div className="flex items-center gap-3">
                      <div className="size-2 rounded-full bg-neutral-600 shrink-0" />
                      <p className="text-sm text-neutral-400">
                        <span className="font-semibold text-neutral-300">Free:</span> {feature.freeLimit}
                      </p>
                    </div>
                    <div className="flex items-center gap-3">
                      <div className="size-2 rounded-full bg-emerald-500 shrink-0 shadow-[0_0_8px_#10b981]" />
                      <p className="text-sm text-emerald-100">
                        <span className="font-semibold text-emerald-400">Premium:</span> {feature.premiumLimit}
                      </p>
                    </div>
                  </div>

                  <Button
                    viewType="confirm"
                    className="rounded-xl px-8 py-3 font-semibold shadow-md shadow-emerald-500/10"
                  >
                    Get Premium
                  </Button>
                </div>

                <div className="flex-1 w-full max-w-md aspect-square relative group">
                  <div
                    className="absolute inset-0 bg-linear-to-br from-neutral-800 to-neutral-900 rounded-[40px] border
                      border-neutral-700/50 shadow-2xl overflow-hidden"
                  >
                    <div
                      className="absolute top-0 right-0 w-full h-full bg-emerald-500/5 group-hover:bg-emerald-500/10
                        transition-colors duration-500"
                    />

                    <div className="absolute inset-0 flex items-center justify-center">
                      <feature.icon
                        className="w-40 h-40 text-neutral-700 group-hover:text-emerald-500/20 group-hover:scale-110
                          transition-all duration-700 ease-out"
                      />
                    </div>

                    <div
                      className="absolute top-6 left-6 w-8 h-8 border-t-4 border-l-4 border-emerald-500/30
                        rounded-tl-xl"
                    />
                    <div
                      className="absolute bottom-6 right-6 w-8 h-8 border-b-4 border-r-4 border-emerald-500/30
                        rounded-br-xl"
                    />
                  </div>
                </div>
              </div>
            )
          })}
        </section>
      </div>
    </div>
  )
}
