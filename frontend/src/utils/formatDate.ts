import { formatDistanceToNow, format, isToday, isYesterday } from 'date-fns'
import { ru } from 'date-fns/locale'

export function formatMessageTime(dateStr: string): string {
  const date = new Date(dateStr)
  return format(date, 'HH:mm')
}

export function formatConversationTime(dateStr: string): string {
  const date = new Date(dateStr)
  if (isToday(date)) return format(date, 'HH:mm')
  if (isYesterday(date)) return 'Вчера'
  return format(date, 'd MMM', { locale: ru })
}

export function formatRelativeTime(dateStr: string): string {
  return formatDistanceToNow(new Date(dateStr), { addSuffix: true, locale: ru })
}
