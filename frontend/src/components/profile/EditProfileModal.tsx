import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '../ui/Input'
import { Button } from '../ui/Button'
import type { User } from '../../types/user'

const schema = z.object({
  username: z.string().min(3).max(30),
  bio: z.string().max(200).optional(),
})

type FormData = z.infer<typeof schema>

interface Props {
  user: User
  onSave: (data: FormData) => Promise<void>
  onClose: () => void
  saving: boolean
}

export function EditProfileModal({ user, onSave, onClose, saving }: Props) {
  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { username: user.username, bio: user.bio },
  })

  return (
    <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4" onClick={onClose}>
      <div className="bg-white rounded-2xl p-6 w-full max-w-sm shadow-xl" onClick={(e) => e.stopPropagation()}>
        <h2 className="text-lg font-semibold mb-4">Редактировать профиль</h2>
        <form onSubmit={handleSubmit(onSave)} className="flex flex-col gap-4">
          <Input label="Имя пользователя" {...register('username')} error={errors.username?.message} />
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-gray-700">О себе</label>
            <textarea
              {...register('bio')}
              placeholder="Расскажите о себе..."
              rows={3}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 resize-none"
            />
          </div>
          <div className="flex gap-2 justify-end pt-2">
            <Button type="button" variant="ghost" onClick={onClose}>Отмена</Button>
            <Button type="submit" loading={saving}>Сохранить</Button>
          </div>
        </form>
      </div>
    </div>
  )
}
