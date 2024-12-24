import { joinURL } from 'ufo'

export default defineEventHandler(async (event) => {
   const config = useRuntimeConfig()
    const proxyUrl = config.proxyApi
    const path = event.path.replace(/^\/api\//, '')

    const target = joinURL(proxyUrl, path)

    return proxyRequest(event, target)
})
