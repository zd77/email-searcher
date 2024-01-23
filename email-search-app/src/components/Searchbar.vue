<template>
    <div class="flex items-center m-4 mb-8">
        <div class="relative flex-grow">
            <span class="absolute top-3 left-3 text-gray-500">🔍</span>
            <input 
                v-model="searchQuery"
                @keyup.enter="handleEnter(searchQuery)"
                type="text"
                class="border border-gray-300 p-3 pl-10 pr-3 rounded focus:outline-none focus:border-blue-500 w-full"
                placeholder="Search..."
            />
        </div>
    </div>
</template>

<script setup>
    import {ref} from 'vue'
    const searchQuery = ref('');
    const { handleEnter } = defineProps(['handleEnter'])
</script>