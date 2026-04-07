import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '../ui/Input'
import { Button } from '../ui/Button'
import type { RegisterRequest } from '../../types/auth'

const schema = z.object({
  email: z.string().email('Неверный email'),
  username: z.string().min(3, 'Минимум 3 символа').max(30, 'Максимум 30 символов'),
  password: z.string().min(8, 'Минимум 8 символов'),
})

interface Props {
  onSubmit: (data: RegisterRequest) => void
  loading: boolean
  error: string | null
}

export function RegisterForm({ onSubmit, loading, error }: Props) {
  const { register, handleSubmit, formState: { errors } } = useForm<RegisterRequest>({
    resolver: zodResolver(schema),
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <Input label="Email" type="email" autoComplete="email" {...register('email')} error={errors.email?.message} />
      <Input label="Имя пользователя" type="text" autoComplete="username" {...register('username')} error={errors.username?.message} />
      <Input label="Пароль" type="password" autoComplete="new-password" {...register('password')} error={errors.password?.message} />
      {error && <p className="text-sm text-red-600 text-center">{error}</p>}
      <Button type="submit" loading={loading} className="w-full">
        Зарегистрироваться
      </Button>
    </form>
  )
}
