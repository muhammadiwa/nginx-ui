import { http } from '@/utils/http'
import { AxiosPromise } from 'axios'

export interface AppSecConfig {
  enabled: boolean
  policies: Policy[]
  rules: Rule[]
}

export interface Policy {
  id: string
  name: string
  description: string
  enabled: boolean
  rules: string[]
}

export interface Rule {
  id: string
  name: string
  description: string
  type: string
  action: 'allow' | 'block' | 'monitor'
  conditions: RuleCondition[]
}

export interface RuleCondition {
  field: string
  operator: string
  value: string | number | boolean
}

export const getAppSecConfig = (): AxiosPromise<AppSecConfig> => {
  return http.get('/api/v1/appsec/config')
}

export const updateAppSecConfig = (config: AppSecConfig): AxiosPromise<void> => {
  return http.put('/api/v1/appsec/config', config)
}

export const getAppSecStats = (): AxiosPromise<{
  totalRequests: number
  blockedRequests: number
  alertedRequests: number
  topAttacks: Array<{type: string, count: number}>
}> => {
  return http.get('/api/v1/appsec/stats')
}

export const getAppSecPolicies = (): AxiosPromise<Policy[]> => {
  return http.get('/api/v1/appsec/policies')
}

export const createAppSecPolicy = (policy: Policy): AxiosPromise<Policy> => {
  return http.post('/api/v1/appsec/policies', policy)
}

export const updateAppSecPolicy = (id: string, policy: Policy): AxiosPromise<Policy> => {
  return http.put(`/api/v1/appsec/policies/${id}`, policy)
}

export const deleteAppSecPolicy = (id: string): AxiosPromise<void> => {
  return http.delete(`/api/v1/appsec/policies/${id}`)
}
