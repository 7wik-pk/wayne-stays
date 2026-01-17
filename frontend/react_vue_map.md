# React to Vue 3 Mapping

As we build Wayne Stays, this document tracks React concepts and their Vue 3 (Composition API) equivalents to help bridge the knowledge gap.

## 1. Components
### React: Functional Components
React components are just JavaScript functions that return JSX.
```tsx
const Hero: React.FC = () => {
  return <div>Hello</div>;
};
```

### Vue 3 Equivalent: Single File Components (SFC)
Vue uses `.vue` files with `<script setup>`, `<template>`, and `<style>`.
```html
<script setup lang="ts">
</script>

<template>
  <div>Hello</div>
</template>
```

---

## 2. Templating (JSX vs Templates)
### React: JSX
React puts HTML directly into JavaScript. You use standard JS for logic.
- **Interpolation**: `{variable}`
- **Class attributes**: `className="foo"` (since `class` is reserved in JS)

### Vue 3 Equivalent: Templates
- **Interpolation**: `{{ variable }}`
- **Class attributes**: `class="foo"`

---

## 3. State Management
### React: `useState` Hook
Used to declare reactive state variables.
```tsx
const [scrolled, setScrolled] = useState(false);
// usage
setScrolled(true);
```

### Vue 3 Equivalent: `ref`
```ts
const scrolled = ref(false);
// usage (note the .value)
scrolled.value = true;
```

---

## 4. Side Effects & Lifecycle
### React: `useEffect` Hook
Handles side effects (API calls, subscriptions, DOM manipulation). It combines `onMounted`, `onUpdated`, and `onUnmounted`.
```tsx
useEffect(() => {
  // onMounted / onUpdated logic
  const handleScroll = () => { ... };
  window.addEventListener('scroll', handleScroll);

  // cleanup (onUnmounted logic)
  return () => {
    window.removeEventListener('scroll', handleScroll);
  };
}, []); // [] dependency array means "only run on mount"
```

### Vue 3 Equivalent: `onMounted` & `onUnmounted`
Vue separates these explicitly.
```ts
onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
```

---

## 5. Loops (Rendering Lists)
### React: `.map()`
Use standard JavaScript array methods.
```tsx
{accommodations.map((acc) => (
  <Card key={acc.title} title={acc.title} />
))}
```

### Vue 3 Equivalent: `v-for`
```html
<Card v-for="acc in accommodations" :key="acc.title" :title="acc.title" />
```

---

## 6. Props
### React
Props are arguments passed to the component function. Use interfaces for TypeScript validation.
```tsx
interface Props {
  title: string;
}
const Card: React.FC<Props> = ({ title }) => { ... }
```

### Vue 3 Equivalent: `defineProps`
```ts
defineProps<{
  title: string
}>();
```

---

## 7. Event Handling
### React: `onClick`, `onChange`, etc.
Event handlers are passed as camelCase props.
```tsx
const handleClick = () => {
    alert('Clicked!');
};

<button onClick={handleClick}>Click Me</button>
```
**Passing arguments**:
```tsx
<button onClick={() => updateFilter('all')}>Filter</button>
```

### Vue 3 Equivalent: `@click` (v-on)
```html
<button @click="handleClick">Click Me</button>
```
**Passing arguments**:
```html
<button @click="updateFilter('all')">Filter</button>
```
