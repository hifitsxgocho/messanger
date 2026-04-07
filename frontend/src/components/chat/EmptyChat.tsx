export function EmptyChat() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center gap-4 text-center px-6">
      <div className="w-20 h-20 bg-blue-100 rounded-full flex items-center justify-center">
        <svg className="w-10 h-10 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z" />
        </svg>
      </div>
      <div>
        <h2 className="font-semibold text-gray-700">Выберите чат</h2>
        <p className="text-sm text-gray-400 mt-1">
          Выберите существующий чат или найдите пользователя
        </p>
      </div>
    </div>
  )
}
