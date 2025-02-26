import type { GetListResponse } from '@/types'
import { http } from '@/utils/http'

interface SecurityLog {
  time: string
  event_severity: string
  asset_name: string
  security_action: string
  appsec_incident_type: string
  source_identifier: string
  source_ip: string
  proxy_ip: string
  http_host: string
  http_method: string
  http_response_code: number
  http_uri_path: string
  matched_location: string
  matched_parameter: string
  matched_sample: string
  event_priority: string
  event_topic: string
  event_name: string
  suggested_remediation: string
}

interface LogsFilter {
  severity?: string;
  startDate?: string;
  endDate?: string;
}

class SecurityLogs {
  protected readonly baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  getImportantEvents(filter?: LogsFilter): Promise<SecurityLog[]> {
    const params = new URLSearchParams()
    if (filter?.severity) params.append('severity', filter.severity)
    if (filter?.startDate) params.append('startDate', filter.startDate)
    if (filter?.endDate) params.append('endDate', filter.endDate)
    
    return http.get(`${this.baseUrl}/important?${params.toString()}`)
  }

  getAllEvents(filter?: LogsFilter): Promise<SecurityLog[]> {
    const params = new URLSearchParams()
    if (filter?.severity) params.append('severity', filter.severity)
    if (filter?.startDate) params.append('startDate', filter.startDate)
    if (filter?.endDate) params.append('endDate', filter.endDate)
    
    return http.get(`${this.baseUrl}/all?${params.toString()}`)
  }

  getNotifications(filter?: LogsFilter): Promise<SecurityLog[]> {
    const params = new URLSearchParams()
    if (filter?.severity) params.append('severity', filter.severity)
    if (filter?.startDate) params.append('startDate', filter.startDate)
    if (filter?.endDate) params.append('endDate', filter.endDate)
    
    return http.get(`${this.baseUrl}/notifications?${params.toString()}`)
  }
}

export default new SecurityLogs('/api/security-logs')
