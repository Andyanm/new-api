import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { getConversationLogs } from '../api'

export function ConversationLogsPanel() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const { data, isLoading } = useQuery({
    queryKey: ['conversation-logs', page],
    queryFn: () => getConversationLogs({ p: page, page_size: 20 }),
  })
  type ConversationLogItem = {
    id: number
    username: string
    token_name: string
    model_name: string
    prompt_text: string
    reply_text: string
  }
  const items: ConversationLogItem[] = data?.data?.items ?? []
  const total = data?.data?.total ?? 0

  return (
    <div className='space-y-3'>
      <div className='text-sm text-muted-foreground'>
        {t('Total records')}: {total}
      </div>
      {isLoading ? (
        <div className='rounded-lg border p-6 text-sm'>{t('Loading...')}</div>
      ) : (
        items.map((item) => (
          <div key={item.id} className='rounded-xl border bg-card p-4 shadow-sm'>
            <div className='mb-2 text-xs text-muted-foreground'>
              #{item.id} · {item.username} · {item.token_name} · {item.model_name}
            </div>
            <div className='grid gap-3 md:grid-cols-2'>
              <div>
                <div className='mb-1 text-xs font-medium'>{t('Prompt')}</div>
                <pre className='max-h-48 overflow-auto rounded-md bg-muted p-2 text-xs whitespace-pre-wrap'>{item.prompt_text}</pre>
              </div>
              <div>
                <div className='mb-1 text-xs font-medium'>{t('Response')}</div>
                <pre className='max-h-48 overflow-auto rounded-md bg-muted p-2 text-xs whitespace-pre-wrap'>{item.reply_text}</pre>
              </div>
            </div>
          </div>
        ))
      )}
      <div className='flex gap-2'>
        <button className='rounded border px-3 py-1 text-sm' onClick={() => setPage((p) => Math.max(1, p - 1))}>
          {t('Previous')}
        </button>
        <button className='rounded border px-3 py-1 text-sm' onClick={() => setPage((p) => p + 1)}>
          {t('Next')}
        </button>
      </div>
    </div>
  )
}
