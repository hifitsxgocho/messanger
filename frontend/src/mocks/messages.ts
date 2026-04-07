import type { Message } from '../types/message'

const msgs = (convId: string, pairs: Array<[string, string, number]>): Message[] =>
  pairs.map(([senderId, body, minutesAgo], i) => ({
    id: `${convId}-msg-${i + 1}`,
    conversationId: convId,
    senderId,
    body,
    createdAt: new Date(Date.now() - minutesAgo * 60 * 1000).toISOString(),
    readAt: senderId !== 'user-1' ? new Date().toISOString() : null,
  }))

export const MOCK_MESSAGES: Record<string, Message[]> = {
  'conv-1': msgs('conv-1', [
    ['user-2', 'Привет! Как дела?', 30],
    ['user-1', 'Всё отлично, работаю над новым проектом', 28],
    ['user-2', 'О, интересно! Расскажи подробнее', 25],
    ['user-1', 'Делаю мессенджер на React + Go', 20],
    ['user-2', 'Звучит круто! Какой стек используешь?', 15],
    ['user-1', 'React 18, TypeScript, Redux, Tailwind на фронте. Go + PostgreSQL на бэке', 10],
    ['user-2', 'Отличный выбор! Tailwind v4 уже попробовал?', 7],
    ['user-1', 'Да, как раз с ним и работаю', 5],
    ['user-2', 'Привет! Как дела?', 2],
  ]),
  'conv-2': msgs('conv-2', [
    ['user-1', 'Привет, bob!', 180],
    ['user-3', 'Привет! Что нового?', 170],
    ['user-1', 'Разбираюсь с pgx для PostgreSQL', 160],
    ['user-3', 'Хороший выбор! Смотри, нашёл крутую библиотеку для Go', 120],
  ]),
  'conv-3': msgs('conv-3', [
    ['user-4', 'Привет! Не подскажешь как сделать sticky header на Tailwind?', 2880],
    ['user-1', 'Конечно! Используй `sticky top-0 z-10`', 2870],
    ['user-4', 'Спасибо за помощь!', 2860],
  ]),
}
