import { Ref } from 'vue'

const VITE_BASE_URL = import.meta.env.VITE_BASE_URL
const SEARCH_MAX_RESULT = 50

export type Items = {
    bodyMsg : string
    subject : string
    from    : string
    date    : string
    to      : string
}

type ApiResponse = {
    hits: {
        total: { value: number };
        hits: Array<{_source: Items}>
    }
}

export const handleSearch = async ( 
    queryString: string, 
    totalEntries: Ref<number>, 
    totalPages: Ref<number>, 
    page = 1 
): Promise<Array<Items>> => {
    try {
        const credentials = btoa('admin:Complexpass#123')
        const response = await fetch(`${VITE_BASE_URL}/api/enron_emails/_search`,{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Basic ${credentials}`,
            },
            body: JSON.stringify({
                "search_type": "match",
                "query": {
                    "term": `${queryString}`,
                    "field": "_all",
                },
                "sort_fields": ["-@Date"],
                "from": (page - 1) * SEARCH_MAX_RESULT,
                "max_results": SEARCH_MAX_RESULT,
                "_source": []
            })
        })
        const formattedResp: ApiResponse = await response.json()
        totalEntries.value = formattedResp.hits.total.value
        totalPages.value = Math.ceil(totalEntries.value/SEARCH_MAX_RESULT)
        const extractData = formattedResp.hits.hits
        return extractData.map(({ _source }: {  _source: Items }) => ({
            to: _source.to,
            from: _source.from,
            subject: _source.subject,
            date: _source.date,
            bodyMsg: _source.bodyMsg,
        }))
    } catch( err ) {
        console.log("Error trying to fetch: ", err)
        return []
    }
}