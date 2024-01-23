<template>
    <header class="text-white p-4 flex items-center">
        <div class="flex items-center">
            <span class="text-2xl mr-2">🚀</span>
            <h1 class="text-xl font-bold">{{ title }}</h1>
        </div>
    </header>  
</template>


<script setup lang='ts'>
    const title = 'Email searcher'
</script>

<style>
header {
  background-color: #000035;
}
</style>