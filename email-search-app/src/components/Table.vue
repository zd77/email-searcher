<template>
    <table class="w-full table-auto">
        <thead>
            <tr>
                <th class="border px-4 py-2 bg-gray-300">Subject</th>
                <th class="border px-4 py-2 bg-gray-300">To</th>
                <th class="border px-4 py-2 bg-gray-300">From</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="(item, index) in data" :key="index" @click="selectRow(item)" class="cursor-pointer hover:bg-gray-100">
                <td class="border px-4 py-2 cursor-pointer max-w-xs overflow-hidden whitespace-nowrap">{{ item.subject }}</td>
                <td class="border px-4 py-2 max-w-xs overflow-hidden whitespace-nowrap">{{ item.to }}</td>
                <td class="border px-4 py-2 max-w-xs overflow-hidden whitespace-nowrap">{{ item.from }}</td>
            </tr>
        </tbody>
    </table>
</template>

<script setup lang="ts">
    const { data, selectRow } = defineProps(['data', 'selectRow'])
</script>