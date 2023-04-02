<template>
  <div class="container">
    <div class="downloads-box">
      <button class="downloads-button">Windows</button>
      <button class="downloads-button">macOS</button>
      <button class="downloads-button">Linux</button>
    </div>
  </div>
</template>

<style>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 90vh;
}

.downloads-box {
  display: flex;
  flex-direction: column;

  height: 50vh;
  width: 50vw;
  padding: 1rem;
  background-color: rgb(0 0 0 / 80%);
  border-radius: 10px;
  justify-content: center;
}

.downloads-button {
  background-color: teal;
  border-radius: 10px;
  margin: 0.5rem 5rem;
}
</style>
