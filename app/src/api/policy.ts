import request from '@/utils/request'

export const getLocalPolicy = () => {
  return request({
    url: '/api/local_policy',
    method: 'get'
  })
}

export const saveLocalPolicy = (content: string) => {
  return request({
    url: '/api/local_policy',
    method: 'post',
    data: { content }
  })
}
