<template>
  <nav>
    <h2>Kilnwood</h2>
    <ul>
      <!-- <li><a href="downloads">Downloads</a></li> -->
    </ul>
  </nav>
</template>

<style>
* {
  padding: 0;
  margin: 0;
}

body {
  font-family: Arial, Tahoma, Serif;
}

nav {
  max-height: 1vh;
  display: flex;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  background: rgb(0 0 0 / 90%);
  border-bottom-right-radius: 5px;
  border-bottom-left-radius: 5px;
  align-items: center;
  color: white;
}

nav ul {
  display: flex;
  list-style: none;
}

nav li {
  padding-left: 1rem;
}

nav a {
  text-decoration: none;
  color: grey;
}
</style>
