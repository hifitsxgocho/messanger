import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '../ui/Input'
import { Button } from '../ui/Button'
import type { LoginRequest } from '../../types/auth'

const schema = z.object({
  email: z.string().email('Неверный email'),
  password: z.string().min(1, 'Введите пароль'),
})

interface Props {
  onSubmit: (data: LoginRequest) => void
  loading: boolean
  error: string | null
}

export function LoginForm({ onSubmit, loading, error }: Props) {
  const { register, handleSubmit, formState: { errors } } = useForm<LoginRequest>({
    resolver: zodResolver(schema),
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <Input label="Email" type="email" autoComplete="email" {...register('email')} error={errors.email?.message} />
      <Input label="Пароль" type="password" autoComplete="current-password" {...register('password')} error={errors.password?.message} />
      {error && <p className="text-sm text-red-600 text-center">{error}</p>}
      <Button type="submit" loading={loading} className="w-full">
        Войти
      </Button>
    </form>
  )
}
