import { createRouter, createWebHistory } from 'vue-router'

import AssetComponent from '@/components/workflows/AssetComponent.vue'
import AssetListComponent from '@/components/workflows/AssetListComponent.vue'

import WorkflowTaskComponent from '@/components/workflows/WorkflowTaskComponent.vue'
import WorkflowTaskListComponent from '@/components/workflows/WorkflowTaskListComponent.vue'

import SchedulerListComponent from '@/components/scheduler/SchedulerListComponent.vue'
import SchedulerComponent from '@/components/scheduler/SchedulerComponent.vue'

import SessionComponent from '@/components/wmpci/SessionComponent.vue'
import SessionListComponent from '@/components/wmpci/SessionListComponent.vue'

import ConnectorComponent from '@/components/wmpci/ConnectorComponent.vue'
import WmpdataComponentDemo from '@/components/wmpdata/WmpdataComponentDemo.vue'


const routes = [

  {
    path: "/",
    name: "home",
    component: WmpdataComponentDemo,
  },

  {
    path: "/wmpci/sessions",
    name: "wmpci session list",
    component: SessionListComponent,
  },
  {
    path: "/wmpci/sessions/:session_id",
    name: "wmpci session detailed",
    component: SessionComponent,
  },

  {
    path: "/wmpci/connector",
    name: "wmpci session connector list",
    component: null,
  },
  {
    path: "/wmpci/connector/:connector_name",
    name: "wmpci session connector config / detailed",
    component: ConnectorComponent,
  },
 
  {
    path: "/scheduler",
    name: "scheduler list",
    component: SchedulerListComponent,
  },

  {
    path: "/scheduler/:scheduler_id",
    name: "scheduler detailed",
    component: SchedulerComponent,
  },

  {
    path: "/asset",
    name: "asset list",
    component: AssetListComponent, 
  },
  {
    path: "/asset/:asset_id",
    name: "asset detailed", 
    component: AssetComponent,
  },

  {
    path: "/workflow/task",
    name: "workflowtask list",
    component: WorkflowTaskListComponent,
  },
  {
    path: "/workflow/task/:task_id",
    name: "workflowtask detailed",
    component: WorkflowTaskComponent,
  }

]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
