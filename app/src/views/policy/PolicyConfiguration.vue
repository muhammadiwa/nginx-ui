<template>
  <div class="policy-container">
    <a-row :gutter="16">
      <a-col :span="12">
        <a-card title="Policy Configuration">
          <template #extra>
            <a-space>
              <a-button type="primary" :loading="loading" @click="fetchPolicy">
                <template #icon>
                  <ReloadOutlined />
                </template>
                Refresh
              </a-button>
              <a-button type="primary" :loading="loading" @click="handleSave">
                <template #icon>
                  <SaveOutlined />
                </template>
                Save
              </a-button>
              <a-button type="primary" :loading="applyLoading" @click="handleApply">
                <template #icon>
                  <CheckOutlined />
                </template>
                Apply
              </a-button>
            </a-space>
          </template>

          <div class="editor-container" :class="{ 'dark-mode': isDarkMode }"
            style="height: 600px; border: 1px solid var(--border-color); border-radius: 4px;">
            <VAceEditor v-model:value="formState.content" lang="yaml" :theme="currentTheme"
              style="height: 100%; font-family: 'JetBrains Mono', monospace; font-size: 14px;" :options="{
                useWorker: false,
                showPrintMargin: true,
                printMarginColumn: 80,
                showGutter: true,
                highlightActiveLine: true,
                enableLiveAutocompletion: true,
                enableSnippets: true,
                displayIndentGuides: true,
                tabSize: 2,
                wrap: false,
                fontSize: '14px',
                scrollPastEnd: 0.5,
                cursorStyle: 'smooth',
                behavioursEnabled: true,
                wrapBehavioursEnabled: true
              }" @init="handleEditorMounted" />
          </div>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Policy Example">
          <a-tabs>
            <a-tab-pane key="policies" tab="Policies">
              <h2>Policies</h2>

              <p class="description">
                The policies section defines the default behaviors and security rules that will be applied across your web infrastructure. It consists of two main parts:
                
                <ul class="description-list">
                  <li>
                    <strong>Default Policy:</strong> Sets the baseline security behavior for all web resources exposed by your NGINX server, Kong Gateway or APISIX Gateway.
                  </li>
                  <li>
                    <strong>Specific Rules:</strong> Allows overriding the default behavior for specific hostname/path combinations, enabling granular security control.
                  </li>
                </ul>
              </p>


              <a-collapse>
                <a-collapse-panel key="1" header="Example">
                  <div class="code-container">
                    <div class="code-actions">
                      <a-tooltip title="Copy to clipboard">
                        <a-button type="text" @click="copyExamplePolicies">
                          <template #icon>
                            <CopyOutlined />
                          </template>
                          Copy
                        </a-button>
                      </a-tooltip>
                    </div>
                    <div class="code-editor">
                      <VAceEditor v-model:value="examplePolicies" lang="yaml" :theme="currentTheme"
                        style="height: 500px; width: 100%" :options="{
                          readOnly: true,
                          showPrintMargin: true,
                          printMarginColumn: 80,
                          showGutter: true,
                          highlightActiveLine: false,
                          useWorker: false,
                          displayIndentGuides: true,
                          fontSize: '14px',
                          fontFamily: 'JetBrains Mono, monospace',
                          showLineNumbers: true,
                          wrap: false,
                        }" />
                    </div>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="2" header="Specification">
                  <a-descriptions bordered>
                    <a-descriptions-item label="default" span={3}>
                      <div class="spec-content">
                        <a-card class="spec-card">
                          <h4><strong>mode</strong> <span class="type">string enum</span></h4>
                          <p>Security engines operation mode. Blocking will only happen in
                            <strong>prevent-learn</strong> mode
                          </p>
                          <a-tag-group>
                            <a-tag color="blue"><strong>prevent-learn</strong></a-tag>
                            <a-tag color="green"><strong>detect-learn</strong></a-tag>
                            <a-tag color="orange"><strong>prevent</strong></a-tag>
                            <a-tag color="cyan"><strong>detect</strong></a-tag>
                            <a-tag><strong>inactive</strong></a-tag>
                          </a-tag-group>
                          <div class="note">Note: <strong>prevent</strong> and <strong>detect</strong> are aliases for
                            <strong>prevent-learn</strong> and <strong>detect-learn</strong>
                          </div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>practices</strong> <span class="type">array of strings</span></h4>
                          <p>Defines which security engines to activate and their specific settings</p>
                          <div class="reference">Reference to <strong>Practice</strong> resource(s)</div>
                          <div class="note">Note: Although this is an array, currently only one practice of each kind
                            can be
                            specified</div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>triggers</strong> <span class="type">array of strings</span></h4>
                          <p>Defines logging verbosity and destination (<strong>stdout</strong>,
                            <strong>syslog</strong>,
                            <strong>cloud</strong>, etc)
                          </p>
                          <div class="reference">Reference to <strong>LogTrigger</strong> resource(s)</div>
                          <div class="note">Note: Although this is an array, currently only a single trigger is
                            supported
                          </div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>custom-response</strong> <span class="type">string</span></h4>
                          <p>Defines prevent mode behaviors upon decision to block:</p>
                          <ul>
                            <li><strong>HTTP response code</strong></li>
                            <li><strong>Block page</strong></li>
                            <li><strong>HTTP redirect</strong></li>
                          </ul>
                          <div class="reference">Reference to <strong>CustomResponse</strong> resource</div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>source-identifiers</strong> <span class="type">string</span></h4>
                          <p>Defines how ML engine will distinguish between sources based on:</p>
                          <ul>
                            <li><strong>IP address</strong></li>
                            <li><strong>X-Forward-For</strong></li>
                            <li><strong>Key in Header/Cookie/JWT</strong></li>
                          </ul>
                          <div class="reference">Reference to <strong>SourcesIdentifier</strong> resource</div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>trusted-sources</strong> <span class="type">string</span></h4>
                          <p>Defines which traffic sources are very unlikely to be malicious. Used for Machine Learning
                            engine.</p>
                          <div class="reference">Reference to <strong>Exception</strong> resource(s)</div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>exceptions</strong> <span class="type">array of strings</span></h4>
                          <p>Defines exceptions to be applied based on:</p>
                          <ul>
                            <li><strong>countryCode</strong></li>
                            <li><strong>countryName</strong></li>
                            <li><strong>sourceIP</strong></li>
                            <li><strong>URL</strong></li>
                            <li><strong>hostName</strong></li>
                            <li><strong>sourceIdentifier</strong></li>
                          </ul>
                          <div class="reference">Reference to <strong>Exception</strong> resource(s)</div>
                        </a-card>
                      </div>
                    </a-descriptions-item>

                    <a-descriptions-item label="specific-rules" span={3}>
                      <div class="spec-content">
                        <a-card class="spec-card">
                          <h4><strong>List of per-host policies</strong></h4>
                          <p>One or more per-host (ingress rule) policies that will override the defaults above</p>

                          <p><strong>host</strong> <span class="type">string</span></p>
                          <p>Policy will apply to this host network path (exactly as appear in ingress rules)</p>

                          <div class="note">All other keys can be used same as described above for
                            <strong>default</strong>.
                          </div>
                        </a-card>
                      </div>
                    </a-descriptions-item>
                  </a-descriptions>
                </a-collapse-panel>
              </a-collapse>
            </a-tab-pane>
            <a-tab-pane key="practices" tab="Practices">
              <h2>Practices</h2>

              <p class="description">
                Practice resources define which security engines will be active and what their settings are.
              </p>

              <a-collapse>
                <a-collapse-panel key="1" header="Example">
                  <div class="code-container">
                    <div class="code-actions">
                      <a-tooltip title="Copy to clipboard">
                        <a-button type="text" @click="copyExamplePractices">
                          <template #icon>
                            <CopyOutlined />
                          </template>
                          Copy
                        </a-button>
                      </a-tooltip>
                    </div>
                    <div class="code-editor">
                      <VAceEditor v-model:value="examplePractices" lang="yaml" :theme="currentTheme"
                        style="height: 500px; width: 100%" :options="{
                          readOnly: true,
                          showPrintMargin: true,
                          printMarginColumn: 80,
                          showGutter: true,
                          highlightActiveLine: false,
                          useWorker: false,
                          displayIndentGuides: true,
                          fontSize: '14px',
                          fontFamily: 'JetBrains Mono, monospace',
                          showLineNumbers: true,
                          wrap: false,
                        }" />
                    </div>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="2" header="Specification">
                  <a-descriptions bordered>
                    <a-descriptions-item label="web-attacks" span={3}>
                      <div class="spec-content">
                        <a-card class="spec-card">
                          <h4><strong>override-mode</strong> <span class="type">string enum</span></h4>
                          <p>Allows overriding the mode defined at the Policy level for this specific engine</p>
                          <a-tag-group>
                            <a-tag color="blue"><strong>prevent-learn</strong></a-tag>
                            <a-tag color="green"><strong>detect-learn</strong></a-tag>
                            <a-tag color="orange"><strong>prevent</strong></a-tag>
                            <a-tag color="cyan"><strong>detect</strong></a-tag>
                            <a-tag><strong>inactive</strong></a-tag>
                            <a-tag color="purple"><strong>as-top-level</strong></a-tag>
                          </a-tag-group>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>minimum-confidence</strong> <span class="type">string enum</span></h4>
                          <p>Defines which security engines to activate and their specific settings</p>
                          <a-tag-group>
                            <a-tag color="orange"><strong>medium</strong></a-tag>
                            <a-tag color="blue"><strong>high</strong></a-tag>
                            <a-tag color="red"><strong>critical</strong></a-tag>
                          </a-tag-group>
                          <div class="note">Default: high</div>
                        </a-card>

                        <a-card class="spec-card">
                          <h4><strong>Size Limits</strong></h4>
                          <div class="size-limits">
                            <div class="limit-item">
                              <p><strong>max-url-size-bytes</strong> <span class="type">integer</span></p>
                              <p>Default: <strong>32768</strong></p>
                            </div>
                            <div class="limit-item">
                              <p><strong>max-object-depth</strong> <span class="type">integer</span></p>
                              <p>Default: <strong>40</strong></p>
                            </div>
                            <div class="limit-item">
                              <p><strong>max-body-size-kb</strong> <span class="type">integer</span></p>
                              <p>Default: <strong>102400</strong></p>
                            </div>
                            <div class="limit-item">
                              <p><strong>max-header-size-bytes</strong> <span class="type">integer</span></p>
                              <p>Default: <strong>32768</strong></p>
                            </div>
                          </div>
                        </a-card>


                        <a-card class="spec-card">
                          <h4><strong>protections</strong></h4>
                          <div class="protection-item">
                            <p><strong>csrf-enabled</strong> <span class="type">string</span></p>
                            <p>Cross Site Request Forgery protection</p>
                            <a-tag-group>
                              <a-tag color="blue">prevent-learn</a-tag>
                              <a-tag color="green">detect-learn</a-tag>
                              <a-tag color="orange">prevent</a-tag>
                              <a-tag color="cyan">detect</a-tag>
                              <a-tag>inactive</a-tag>
                            </a-tag-group>
                            <div class="note">Default: inactive</div>
                          </div>

                          <div class="protection-item">
                            <p><strong>error-disclosure-enabled</strong> <span class="type">string</span></p>
                            <p>Prevent disclosure of technical information to the attacker in server error messages</p>
                            <a-tag-group>
                              <a-tag color="blue">prevent-learn</a-tag>
                              <a-tag color="green">detect-learn</a-tag>
                              <a-tag color="orange">prevent</a-tag>
                              <a-tag color="cyan">detect</a-tag>
                              <a-tag>inactive</a-tag>
                            </a-tag-group>
                            <div class="note">Default: inactive</div>
                          </div>

                          <div class="protection-item">
                            <p><strong>open-redirect-enabled</strong> <span class="type">string</span></p>
                            <p>Protect against URL redirection to untrusted sites</p>
                            <a-tag-group>
                              <a-tag color="blue">prevent-learn</a-tag>
                              <a-tag color="green">detect-learn</a-tag>
                              <a-tag color="orange">prevent</a-tag>
                              <a-tag color="cyan">detect</a-tag>
                              <a-tag>inactive</a-tag>
                            </a-tag-group>
                            <div class="note">Default: inactive</div>
                          </div>

                          <div class="protection-item">
                            <p><strong>non-valid-http-methods</strong> <span class="type">boolean</span></p>
                            <p>Prevent attacker from sending requests with unsafe HTTP methods</p>
                            <a-tag-group>
                              <a-tag color="green">true</a-tag>
                              <a-tag>false</a-tag>
                            </a-tag-group>
                            <div class="note">Default: false</div>
                          </div>
                        </a-card>
                      </div>
                    </a-descriptions-item>

                    <!-- Additional sections for open-api-schema-validation, anti-bot, and snort-signatures -->
                    <a-descriptions-item label="Coming Soon Features" span={3}>
                      <div class="spec-content">
                        <a-card class="spec-card coming-soon">
                          <h4><strong>open-api-schema-validation</strong></h4>
                          <p class="coming-soon-text">Currently not supported yet, will be added soon</p>
                          <!-- Add specifications -->
                        </a-card>

                        <a-card class="spec-card coming-soon">
                          <h4><strong>anti-bot</strong></h4>
                          <p class="coming-soon-text">Currently not supported yet, will be added soon</p>
                          <!-- Add specifications -->
                        </a-card>

                        <a-card class="spec-card coming-soon">
                          <h4><strong>snort-signatures</strong></h4>
                          <p class="coming-soon-text">Currently not supported yet, will be added soon</p>
                          <!-- Add specifications -->
                        </a-card>
                      </div>
                    </a-descriptions-item>
                  </a-descriptions>
                </a-collapse-panel>
              </a-collapse>
            </a-tab-pane>

            <a-tab-pane key="custom-response" tab="Custom Response">
              <h2>Custom Response</h2>

              <p class="description">
                Define how the engine should respond when blocking requests.
              </p>

              <a-collapse>
                <a-collapse-panel key="1" header="Example">
                  <div class="code-container">
                    <div class="code-actions">
                      <a-tooltip title="Copy to clipboard">
                        <a-button type="text" @click="copyCustomResponseExample">
                          <template #icon>
                            <CopyOutlined />
                          </template>
                          Copy
                        </a-button>
                      </a-tooltip>
                    </div>
                    <div class="code-editor">
                      <VAceEditor v-model:value="customResponseExample" lang="yaml" :theme="currentTheme"
                        style="height: 300px; width: 100%" :options="{
                          readOnly: true,
                          showPrintMargin: true,
                          printMarginColumn: 80,
                          showGutter: true,
                          highlightActiveLine: false,
                          useWorker: false,
                          displayIndentGuides: true,
                          fontSize: '14px',
                          fontFamily: 'JetBrains Mono, monospace',
                          showLineNumbers: true,
                          wrap: false,
                        }" />
                    </div>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="2" header="Specification">
                  <div class="spec-content">
                    <a-card class="spec-card">
                      <h4><strong>mode</strong> <span class="type">enum</span></h4>
                      <p>Engine will take one of these actions upon decision to block request</p>
                      <a-tag-group>
                        <a-tag color="blue"><strong>block-page</strong></a-tag>
                        <a-tag color="green"><strong>response-code-only</strong></a-tag>
                      </a-tag-group>
                      <ul>
                        <li><strong>block-page</strong>: send HTML with text to client + HTTP response code</li>
                        <li><strong>response-code-only</strong>: send only response code</li>
                      </ul>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>message-title</strong> <span class="type">string</span></h4>
                      <p>Title of block page that will be displayed only in case mode is block page and engine decided
                        to
                        block</p>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>message-body</strong> <span class="type">string</span></h4>
                      <p>Content of block page that will be displayed only in case mode is block page and engine decided
                        to
                        block</p>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>http-response-code</strong> <span class="type">integer</span></h4>
                      <p>HTTP code that will be returned to client upon engine decision to block</p>
                      <div class="note">
                        <strong>Range:</strong> 100-599<br />
                        <strong>Default:</strong> 403 (HTTP Forbidden)
                      </div>
                    </a-card>
                  </div>
                </a-collapse-panel>
              </a-collapse>
            </a-tab-pane>

            <a-tab-pane key="log-trigger" tab="Log Trigger">
              <h2>Log Trigger</h2>

              <p class="description">
                Configure logging settings and destinations for various security events.
              </p>

              <a-collapse>
                <a-collapse-panel key="1" header="Example">
                  <div class="code-container">
                    <div class="code-actions">
                      <a-tooltip title="Copy to clipboard">
                        <a-button type="text" @click="copyLogTriggerExample">
                          <template #icon>
                            <CopyOutlined />
                          </template>
                          Copy
                        </a-button>
                      </a-tooltip>
                    </div>
                    <div class="code-editor">
                      <VAceEditor v-model:value="logTriggerExample" lang="yaml" :theme="currentTheme"
                        style="height: 500px; width: 100%" :options="{
                          readOnly: true,
                          showPrintMargin: true,
                          printMarginColumn: 80,
                          showGutter: true,
                          highlightActiveLine: false,
                          useWorker: false,
                          displayIndentGuides: true,
                          fontSize: '14px',
                          fontFamily: 'JetBrains Mono, monospace',
                          showLineNumbers: true,
                          wrap: false,
                        }" />
                    </div>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="2" header="Specification">
                  <div class="spec-content">
                    <a-card class="spec-card">
                      <h4><strong>access-control-logging</strong></h4>
                      <p>Configure logging for Access Control events</p>

                      <div class="config-item">
                        <p><strong>allow-events</strong> <span class="type">boolean</span></p>
                        <p>Log access control allow events</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: false</div>
                      </div>

                      <div class="config-item">
                        <p><strong>drop-events</strong> <span class="type">boolean</span></p>
                        <p>Log access control drop events</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>additional-suspicious-events-logging</strong></h4>
                      <p>Configure additional logging for suspicious events based on a selectable minimum severity-level
                      </p>

                      <div class="config-item">
                        <p><strong>enabled</strong> <span class="type">boolean</span></p>
                        <p>Enable/disable additional suspicious events logging</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>

                      <div class="config-item">
                        <p><strong>minimum-severity</strong> <span class="type">string enum</span></p>
                        <p>Select minimum severity level</p>
                        <a-tag-group>
                          <a-tag color="blue">high</a-tag>
                          <a-tag color="red">critical</a-tag>
                        </a-tag-group>
                        <div class="note">Default: high</div>
                      </div>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>waf-logging</strong></h4>
                      <p>Configure logging for open-waf events (threat prevention, machine learning)</p>

                      <div class="config-item">
                        <p><strong>detect-events</strong> <span class="type">boolean</span></p>
                        <p>Log detected events</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>

                      <div class="config-item">
                        <p><strong>prevent-events</strong> <span class="type">boolean</span></p>
                        <p>Log prevented events</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>

                      <div class="config-item">
                        <p><strong>all-web-requests</strong> <span class="type">boolean</span></p>
                        <p>Log all web requests (has performance impact!)</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: false</div>
                      </div>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>extended-logging</strong></h4>

                      <div class="config-item">
                        <p><strong>url-path</strong> <span class="type">boolean</span></p>
                        <p>Log URL path</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>

                      <div class="config-item">
                        <p><strong>url-query</strong> <span class="type">boolean</span></p>
                        <p>Log URL query</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: true</div>
                      </div>

                      <div class="config-item">
                        <p><strong>http-headers</strong> <span class="type">boolean</span></p>
                        <p>Log the HTTP headers (has performance impact!)</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: false</div>
                      </div>

                      <div class="config-item">
                        <p><strong>request-body</strong> <span class="type">boolean</span></p>
                        <p>Log the request body (has performance impact)</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: false</div>
                      </div>
                    </a-card>

                    <a-card class="spec-card">
                      <h4><strong>log-destination</strong></h4>

                      <div class="config-item">
                        <p><strong>cloud</strong> <span class="type">boolean</span></p>
                        <p>Enable or disable logging to the waf-open Cloud Service</p>
                        <a-tag-group>
                          <a-tag color="green">true</a-tag>
                          <a-tag>false</a-tag>
                        </a-tag-group>
                        <div class="note">Default: false</div>
                      </div>

                      <div class="config-item">
                        <p><strong>file</strong> <span class="type">string</span></p>
                        <p>Define file path to save logs to</p>
                        <div class="note">Local path from root directory of the open-waf container</div>
                      </div>

                      <div class="config-item">
                        <p><strong>stdout</strong></p>
                        <p>Configure logging to standard-out</p>
                        <div class="sub-item">
                          <h6><strong>format</strong> <span class="type">string enum</span></h6>
                          <a-tag-group>
                            <a-tag color="blue">json</a-tag>
                            <a-tag color="cyan">json-formatted</a-tag>
                          </a-tag-group>
                        </div>
                      </div>

                      <div class="config-item">
                        <p><strong>syslog-service</strong> <span class="type">objects array</span></p>
                        <p>Define one or more syslog servers and corresponding ports</p>
                        <div class="sub-item">
                          <h6><strong>address</strong> <span class="type">string</span></h6>
                          <p>Syslog server IP address</p>
                        </div>
                        <div class="sub-item">
                          <h6><strong>port</strong> <span class="type">integer</span></h6>
                          <p>Syslog server port</p>
                        </div>
                      </div>

                      <div class="config-item">
                        <p><strong>cef-service</strong></p>
                        <p>Allows sending files to a log destination in CEF format</p>
                        <div class="sub-item">
                          <h6><strong>address</strong> <span class="type">string</span></h6>
                          <p>CEF server IP address</p>
                        </div>
                        <div class="sub-item">
                          <h6><strong>port</strong> <span class="type">integer</span></h6>
                          <p>CEF server port</p>
                        </div>
                        <div class="sub-item">
                          <h6><strong>proto</strong> <span class="type">string enum</span></h6>
                          <p>Select the correct protocol</p>
                          <a-tag-group>
                            <a-tag color="blue">tcp</a-tag>
                            <a-tag color="purple">udp</a-tag>
                          </a-tag-group>
                        </div>
                      </div>
                    </a-card>
                  </div>
                </a-collapse-panel>
              </a-collapse>
            </a-tab-pane>

            <a-tab-pane key="exceptions" tab="Exceptions">
  <h2>Exceptions</h2>
  
  <p class="description">
    Define custom exceptions and rules with flexible matching parameters and actions.
  </p>

  <a-collapse>
    <a-collapse-panel key="1" header="Example">
      <div class="code-container">
        <div class="code-actions">
          <a-tooltip title="Copy to clipboard">
            <a-button type="text" @click="copyExceptionsExample">
              <template #icon><CopyOutlined /></template>
              Copy
            </a-button>
          </a-tooltip>
        </div>
        <div class="code-editor">
          <VAceEditor
            v-model:value="exceptionsExample"
            lang="yaml"
            :theme="currentTheme"
            style="height: 600px; width: 100%"
            :options="{
              readOnly: true,
              showPrintMargin: true,
              printMarginColumn: 80,
              showGutter: true,
              highlightActiveLine: false,
              useWorker: false,
              displayIndentGuides: true,
              fontSize: '14px',
              fontFamily: 'JetBrains Mono, monospace',
              showLineNumbers: true,
              wrap: false,
            }"
          />
        </div>
      </div>
    </a-collapse-panel>

    <a-collapse-panel key="2" header="Specification">
      <div class="spec-content">
        <a-card class="spec-card">
          <h4><strong>action</strong> <span class="type">string enum</span></h4>
          <p>Action to be performed when exception matches</p>
          <a-tag-group>
            <a-tag color="blue"><strong>skip</strong></a-tag>
            <a-tag color="green"><strong>accept</strong></a-tag>
            <a-tag color="red"><strong>drop</strong></a-tag>
            <a-tag color="orange"><strong>suppressLog</strong></a-tag>
          </a-tag-group>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>Match Parameters</strong></h4>
          <div class="param-grid">
            <div class="param-item">
              <p><strong>sourceIp</strong> <span class="type">string array</span></p>
              <p>Source IP(s)</p>
            </div>
            <div class="param-item">
              <p><strong>url</strong> <span class="type">string array</span></p>
              <p>URL(s)</p>
            </div>
            <div class="param-item">
              <p><strong>sourceIdentifier</strong> <span class="type">string array</span></p>
              <p>Identified source(s)</p>
            </div>
            <div class="param-item">
              <p><strong>protectionName</strong> <span class="type">string array</span></p>
              <p>Protection(s)</p>
            </div>
            <div class="param-item">
              <p><strong>paramValue</strong> <span class="type">string array</span></p>
              <p>Parameter value(s)</p>
            </div>
            <div class="param-item">
              <p><strong>paramName</strong> <span class="type">string array</span></p>
              <p>Parameter name(s)</p>
            </div>
            <div class="param-item">
              <p><strong>hostName</strong> <span class="type">string array</span></p>
              <p>Host name(s)</p>
            </div>
            <div class="param-item">
              <p><strong>countryCode</strong> <span class="type">string array</span></p>
              <p>Country code(s)</p>
            </div>
            <div class="param-item">
              <p><strong>countryName</strong> <span class="type">string array</span></p>
              <p>Country name(s)</p>
            </div>
          </div>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>comment</strong> <span class="type">string</span></h4>
          <p>Comment for the exception</p>
        </a-card>
      </div>
    </a-collapse-panel>
  </a-collapse>
</a-tab-pane>

<a-tab-pane key="trusted-sources" tab="Trusted Sources">
  <h2>Trusted Sources</h2>
  
  <p class="description">
    Define trusted sources and behavioral patterns for the ML engine to learn benign traffic patterns.
  </p>

  <a-collapse>
    <a-collapse-panel key="1" header="Example">
      <div class="code-container">
        <div class="code-actions">
          <a-tooltip title="Copy to clipboard">
            <a-button type="text" @click="copyTrustedSourcesExample">
              <template #icon><CopyOutlined /></template>
              Copy
            </a-button>
          </a-tooltip>
        </div>
        <div class="code-editor">
          <VAceEditor
            v-model:value="trustedSourcesExample"
            lang="yaml"
            :theme="currentTheme"
            style="height: 200px; width: 100%"
            :options="{
              readOnly: true,
              showPrintMargin: true,
              printMarginColumn: 80,
              showGutter: true,
              highlightActiveLine: false,
              useWorker: false,
              displayIndentGuides: true,
              fontSize: '14px',
              fontFamily: 'JetBrains Mono, monospace',
              showLineNumbers: true,
              wrap: false,
            }"
          />
        </div>
      </div>
    </a-collapse-panel>

    <a-collapse-panel key="2" header="Specification">
      <div class="spec-content">
        <a-card class="spec-card">
          <p class="description">
            Define trusted sources by referencing the source identifiers custom resources as well as setting the minimum amount of sources that need to be observed by the behavioural ML engine sending certain identical traffic patterns in order to learn this behaviour as being benign.
          </p>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>minNumOfSources</strong> <span class="type">integer</span></h4>
          <p>Minimum amount of sources having to be observed sending same traffic patterns to learn behaviour as benign.</p>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>sourcesIdentifiers</strong> <span class="type">string array</span></h4>
          <p>Specify one or more source identifiers</p>
        </a-card>
      </div>
    </a-collapse-panel>
  </a-collapse>
</a-tab-pane>

<a-tab-pane key="source-identifiers" tab="Source Identifiers">
  <h2>Source Identifiers</h2>
  
  <p class="description">
    Define source identifiers that can be used in trusted sources custom resources.
  </p>

  <a-collapse>
    <a-collapse-panel key="1" header="Example">
      <div class="code-container">
        <div class="code-actions">
          <a-tooltip title="Copy to clipboard">
            <a-button type="text" @click="copySourceIdentifiersExample">
              <template #icon><CopyOutlined /></template>
              Copy
            </a-button>
          </a-tooltip>
        </div>
        <div class="code-editor">
          <VAceEditor
            v-model:value="sourceIdentifiersExample"
            lang="yaml"
            :theme="currentTheme"
            style="height: 150px; width: 100%"
            :options="{
              readOnly: true,
              showPrintMargin: true,
              printMarginColumn: 80,
              showGutter: true,
              highlightActiveLine: false,
              useWorker: false,
              displayIndentGuides: true,
              fontSize: '14px',
              fontFamily: 'JetBrains Mono, monospace',
              showLineNumbers: true,
              wrap: false,
            }"
          />
        </div>
      </div>
    </a-collapse-panel>

    <a-collapse-panel key="2" header="Specification">
      <div class="spec-content">
        <a-card class="spec-card">
          <h4><strong>identifiers</strong> <span class="type">objects array</span></h4>
          <p>Provide single source identifier</p>
          <div class="note">
            Note: Although this is an array, currently only adding single source identifier is supported. Future updates will allow multiple source-identifiers for fallback checking.
          </div>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>source-identifier</strong> <span class="type">string enum</span></h4>
          <p>Specify the source identifier type of which the content shall be matched</p>
          <a-tag-group>
            <a-tag color="blue">headerkey</a-tag>
            <a-tag color="green">JWTKey</a-tag>
            <a-tag color="orange">cookie</a-tag>
            <a-tag color="cyan">sourceip</a-tag>
            <a-tag color="purple">x-forwarded-for</a-tag>
          </a-tag-group>
        </a-card>

        <a-card class="spec-card">
          <h4><strong>value</strong> <span class="type">string array</span></h4>
          <p>Content to match the specified sourceIdentifier type:</p>
          <ul class="value-list">
            <li>For types <strong>headerkey</strong>, <strong>cookie</strong> and <strong>JWTKey</strong>: Provide the fieldname that designates user</li>
            <li>For type <strong>Source IP</strong>: No value is required</li>
            <li>For type <strong>x-forwarded-for</strong>: Provide previous proxy hops if there are any</li>
          </ul>
        </a-card>
      </div>
    </a-collapse-panel>
  </a-collapse>
</a-tab-pane>

          </a-tabs>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, SaveOutlined, CheckOutlined, CopyOutlined } from '@ant-design/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import { useSettingsStore } from '@/pinia'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-monokai'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import axios from 'axios'


const formState = ref({
  content: ''
})

const loading = ref(false)
const applyLoading = ref(false)
let editor: any = null

const settings = useSettingsStore()
const isDarkMode = computed(() => settings.theme === 'dark')
const currentTheme = computed(() => isDarkMode.value ? 'monokai' : 'chrome')

const handleEditorMounted = (editorInstance: any) => {
  editor = editorInstance
  editor.renderer.setScrollMargin(10, 10, 10, 10)
  editor.setOptions({
    enableBasicAutocompletion: true,
    enableLiveAutocompletion: true,
    enableSnippets: true
  })

  // Update theme when dark mode changes
  watch(isDarkMode, (newValue) => {
    editor.setTheme(`ace/theme/${newValue ? 'monokai' : 'chrome'}`)
  })
}

const fetchPolicy = async () => {
  try {
    loading.value = true
    const response = await axios.get('/api/policy')
    formState.value.content = response.data.content
  } catch (error) {
    message.error('Failed to fetch policy')
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  const content = formState.value.content.trim()
  if (!content) {
    message.error('Policy content cannot be empty')
    return
  }

  try {
    loading.value = true
    await axios.post('/api/policy', {
      content: content
    })
    message.success('Policy saved successfully')
  } catch (error) {
    message.error('Failed to save policy')
  } finally {
    loading.value = false
  }
}

const handleApply = async () => {
  applyLoading.value = true
  try {
    const response = await axios.post('/api/policy/apply', {}, {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    })

    // Check if output contains "New policy applied"
    if (response.data.output && response.data.output.includes('New policy applied')) {
      message.success('Policy applied successfully')
    } else if (response.data.output && response.data.output.includes("Policy didn't change")) {
      message.warning('Policy did not change. Please verify that you have a valid new policy.')
    } else {
      // Show error message for other cases
      const errorMsg = response.data.error || response.data.output || 'Unknown error occurred'
      message.error('Error Applying Policy: ' + errorMsg)
    }
  } catch (error: any) {
    const errorMsg = error.response?.data?.error || error.response?.data?.output || error.response?.data?.message || error.message || 'Network error'
    message.error('Failed to apply policy: ' + errorMsg)
  } finally {
    applyLoading.value = false
  }
}

const examplePolicies = ref(`policies:
  default:
    triggers:
    - waf-special-log-trigger
    mode: detect-learn
    practices:
    - webapp-default-practice
    source-identifiers: waf-source-identifiers-sourceip-example
    trusted-sources: waf-trusted-source-example
    custom-response: waf-web-user-response-example
    exceptions:
    - waf-exception-example

  specific-rules:
  - host: web.server.com/example
    triggers:
    - waf-special-log-trigger
    mode: prevent-learn
    practices:
    - webapp-best-practice
    source-identifiers: waf-source-identifiers-sourceip-example
    trusted-sources: waf-trusted-source-example
    custom-response: waf-web-user-response-example
    exceptions:
    - waf-exception-example

  - host: web.server.com/another-example
    triggers:
    - waf-special-log-trigger
    mode: prevent-learn
    practices:
    - webapp-best-practice
    source-identifiers: waf-source-identifiers-sourceip-example
    trusted-sources: waf-trusted-source-example
    custom-response: waf-web-user-response-example
    exceptions:
    - waf-exception-example`)

const copyExamplePolicies = () => {
  navigator.clipboard.writeText(examplePolicies.value)
  message.success('Code copied to clipboard')
}

const examplePractices = ref(`practices:
  - name: webapp-best-practice
    openapi-schema-validation:
      files: []
      override-mode: 'prevent'
    snort-signatures:
      files: []
      override-mode: 'prevent'
    web-attacks:
      max-body-size-kb: 1222
      max-header-size-bytes: 44343
      max-object-depth: 2111
      max-url-size-bytes: 34434
      minimum-confidence: high
      override-mode: 'prevent'
      protections:
        csrf-enabled: prevent
        error-disclosure-enabled: prevent
        non-valid-http-methods: true
        open-redirect-enabled: prevent
    anti-bot:
      injected-URIs: []
      validated-URIs: []
      override-mode: 'prevent'`)

const copyExamplePractices = () => {
  navigator.clipboard.writeText(examplePractices.value)
  message.success('Code copied to clipboard')
}

const customResponseExample = ref(`custom-responses:
  - name: waf-default-web-user-response
    mode: response-code-only
    http-response-code: 403

  - name: waf-web-user-response-example
    mode: block-page
    http-response-code: 403
    message-title: Block page title
    message-body: "<h1>Access blocked by open-waf.</h1><p>Your access will be logged.</p>"`)

const copyCustomResponseExample = () => {
  navigator.clipboard.writeText(customResponseExample.value)
  message.success('Code copied to clipboard')
}

const logTriggerExample = ref(`logtriggers:
  - name: waf-special-log-trigger
    access-control-logging:
      allow-events: false
      drop-events: true
    additional-suspicious-events-logging:
      enabled: true
      minimum-severity: high
      response-body: false
    waf-logging:
      all-web-requests: false
      detect-events: true
      prevent-events: true
    extended-logging:
      http-headers: false
      request-body: false
      url-path: false
      url-query: false
    log-destination:
      cloud: false
      file: "/a/b/c"
      stdout:
        format: json
      syslog-service:
      - address: 1.2.3.4
        port: 514
      cef-service:
        address: 5.6.7.8
        port: 514
        proto: tcp`)

const copyLogTriggerExample = () => {
  navigator.clipboard.writeText(logTriggerExample.value)
  message.success('Code copied to clipboard')
}

const exceptionsExample = ref(`exceptions:
  - name: waf-exception-example
    action: skip
    comment: This is an example exception comment
    countryCode:
    - CA
    - IL
    countryName:
    - Israel
    - Canada
    hostName:
    - fff
    paramName:
    - key
    paramValue:
    - rrr
    protectionName:
    - cveee
    sourceIdentifier:
    - david
    sourceIp:
    - 1.2.3.4
    - '3,3,3,3'
    url:
    - "/rrr"

  - name: exception-example-2
    action: accept
    hostName:
    - fff
    url:
    - "/rrr"

  - name: exception-example-3
    action: drop
    comment: This is an example exception comment
    countryName:
    - Israel
    - Canada
    protectionName:
    - cveee
    sourceIdentifier:
    - david
    sourceIp:
    - 1.2.3.4
    - 2.3.4.5
    url:
    - "/rrr"

  - name: exception-example-4
    action: suppressLog
    comment: This is an example exception comment
    countryCode:
    - CA
    - IL
    countryName:
    - Israel
    - Canada
    hostName:
    - fff
    url:
    - "/rrr"`)

const copyExceptionsExample = () => {
  navigator.clipboard.writeText(exceptionsExample.value)
  message.success('Code copied to clipboard')
}

const trustedSourcesExample = ref(`trustedsources:
  - name: waf-trusted-source-example
    minNumOfSources: 3
    sourcesIdentifiers: [0.0.0.0, 1.1.1.1, 2.2.2.2]`)

const copyTrustedSourcesExample = () => {
  navigator.clipboard.writeText(trustedSourcesExample.value)
  message.success('Code copied to clipboard')
}

const sourceIdentifiersExample = ref(`source-identifiers:
  - name: waf-source-identifiers-sourceip-example
    identifiers:
    - sourceIdentifier: sourceip`)

const copySourceIdentifiersExample = () => {
  navigator.clipboard.writeText(sourceIdentifiersExample.value)
  message.success('Code copied to clipboard')
}


onMounted(() => {
  fetchPolicy()
})
</script>

<style scoped>
.policy-container {
  margin: 0 auto;
  max-width: 100%;
}

.description {
  margin: 16px 0;
  color: var(--text-color);
}

.code-block {
  position: relative;
  background: var(--code-bg);
  border-radius: 4px;
  padding: 16px;
}

.code-container {
  position: relative;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}

.code-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.9);
  padding: 4px;
  border-radius: 4px;
}


pre {
  margin: 0;
  padding: 0;
  overflow-x: auto;
}

code {
  font-family: 'JetBrains Mono', monospace;
  white-space: pre;
}

.spec-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.spec-card {
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 12px;
}

.type {
  font-size: 12px;
  color: #1890ff;
  background: #e6f7ff;
  padding: 2px 8px;
  border-radius: 4px;
  margin-left: 8px;
}

.reference {
  font-style: italic;
  color: #52c41a;
  margin-top: 8px;
}

.note {
  font-size: 12px;
  color: #faad14;
  background: #fffbe6;
  padding: 8px;
  border-radius: 4px;
  margin-top: 8px;
}

ul {
  padding-left: 20px;
  margin: 8px 0;
}

:deep(.ant-card) {
  height: 100%;
}

/* CSS Variables for theming */
:root {
  --border-color: #d9d9d9;
  --background-color-light: #f5f5f5;
}

:root[data-theme='dark'] {
  --border-color: #434343;
  --background-color-light: #262626;
}


.editor-container.dark-mode {
  background-color: #272822;
}

.editor-container.dark-mode :deep(.ace_gutter) {
  background-color: #272822;
}

.dark .policy-example {
  background-color: #262626;
  color: #fff;
}

.protection-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
}

.protection-item:last-child {
  border-bottom: none;
}

.coming-soon {
  opacity: 0.7;
}

.coming-soon-text {
  color: #ff4d4f;
  font-style: italic;
}

.spec-card :deep(.ant-descriptions) {
  margin: 16px 0;
}

.spec-card :deep(.ant-descriptions-item-label) {
  font-weight: bold;
  color: var(--heading-color);
}

.size-limits {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.limit-item {
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--background-color-light);
}

.limit-item p {
  font-size: 14px;
  margin: 0 0 8px 0;
}

.limit-item p {
  font-size: 14px;
  margin: 0;
}

/* Ensure consistent font sizes across all spec cards */
.spec-card h4 {
  font-size: 14px;
  margin: 0 0 16px 0;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.spec-card p {
  font-size: 14px;
  margin: 4px 0;
  line-height: 1.4;
}

.config-item {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
}

.config-item:last-child {
  border-bottom: none;
}

.sub-item {
  margin-left: 16px;
  padding: 8px 0;
}

.sub-item h6 {
  font-size: 14px;
  margin: 0 0 8px 0;
}

.sub-item p {
  margin: 4px 0;
}

.param-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
  padding: 16px;
}

.param-item {
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--background-color-light);
  line-height: 1.4;
}

.param-item p {
  font-size: 14px;
  margin: 0 0 8px 0;
}

.param-item p {
  font-size: 14px;
  margin: 0;
  color: var(--text-secondary-color);
}

.value-list {
  list-style-type: none;
  padding: 0;
  margin: 8px 0;
}

.value-list li {
  padding: 6px 0;
  line-height: 1.4;
  border-bottom: 1px solid var(--border-color);
}

.value-list li:last-child {
  border-bottom: none;
}

.description-list {
  margin: 8px 0;
  padding-left: 20px;
}

.description-list li {
  margin: 8px 0;
  line-height: 1.4;
}
</style>
