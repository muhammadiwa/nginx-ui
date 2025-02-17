import { http } from '@/utils/http'

class PolicyApi {
  getPolicy() {
    return http.get('/api/settings/policy')
  }

  savePolicy(content: string) {
    return http.post('/api/settings/policy', { content })
  }
}

export default new PolicyApi()
