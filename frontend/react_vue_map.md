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

---

## 8. Conditional Rendering
### React: `&&` and the ternary operator
JSX has no directives, so you use plain JavaScript expressions. `cond && <X />` renders `<X />` only when `cond` is truthy; `a ? <X /> : <Y />` picks one of two.
```tsx
{loading && <p>Loading accommodations...</p>}

{loading ? (
  <div>Loading collection...</div>
) : error ? (
  <div>Unable to load full collection.</div>
) : (
  <div className="merch-grid">...</div>
)}
```
**Pitfall**: `{items.length && <List />}` renders a literal `0` when the list is empty, because `0` is falsy but still renderable. Use `{items.length > 0 && <List />}`.

### Vue 3 Equivalent: `v-if` / `v-else-if` / `v-else`
```html
<div v-if="loading">Loading collection...</div>
<div v-else-if="error">Unable to load full collection.</div>
<div v-else class="merch-grid">...</div>
```
There is no `v-show` equivalent; to hide without unmounting, toggle a class or `style={{ display: 'none' }}`.

---

## 9. Dynamic Classes & Inline Styles
### React: strings and style objects
`className` is just a string, so dynamic classes are built with template literals. `style` takes an object with camelCase property names.
```tsx
<nav className={`navbar ${scrolled ? 'scrolled' : ''}`}>
<button className={`filter-btn ${filter === cat ? 'active' : ''}`}>
<div style={{ maxWidth: '600px', marginTop: '1rem' }}>
```
In this project we prefer CSS classes over inline styles (each component imports its own `.css` file).

### Vue 3 Equivalent: `:class` and `:style` bindings
Vue accepts objects/arrays and merges them with the static `class` for you.
```html
<nav class="navbar" :class="{ scrolled }">
<button class="filter-btn" :class="{ active: filter === cat }">
<div :style="{ maxWidth: '600px', marginTop: '1rem' }">
```

---

## 10. Fragments
### React: `<>...</>`
A component must return a single element. A fragment groups siblings without adding a wrapper `<div>` to the DOM.
```tsx
// pages/HomePage.tsx
return (
  <>
    <Hero />
    <main id="accommodations">...</main>
    <MerchShop />
  </>
);
```
Inside `.map()`, use the long form so you can give it a key: `<Fragment key={id}>...</Fragment>`.

### Vue 3 Equivalent: multi-root templates
Vue 3 templates can have several root nodes, so you just write them. `<template>` (without `v-if`/`v-for`) is the closest thing to a fragment inside a template.

---

## 11. Effect Dependencies & Cleanup (`useEffect` in depth)
### React: the dependency array decides when an effect re-runs
| Dependency array | Runs |
|---|---|
| `useEffect(fn, [])` | once after the first render (like `onMounted`) |
| `useEffect(fn, [a, b])` | after first render, and again whenever `a` or `b` changes |
| `useEffect(fn)` | after every render (rarely what you want) |

The function returned from the effect is the **cleanup**. It runs before the effect re-runs and when the component unmounts. We use it to cancel in-flight requests:
```ts
// hooks/useApi.ts
useEffect(() => {
  const controller = new AbortController();
  fetcher(controller.signal).then(...);
  return () => controller.abort(); // cleanup
}, [fetcher]);
```
React compares dependencies by identity (`Object.is`). A function or object created during render is "new" every time, which would re-run the effect on every render. That's why `useApi` is given module-level functions like `api.listMerchandise`.

**StrictMode** (see 12) deliberately mounts, unmounts and re-mounts every component once in development, so effects run twice. Cleanups like the `abort()` above make that harmless.

### Vue 3 Equivalent: `watch` / `watchEffect` with `onCleanup`
An effect with dependencies is closer to `watch` than to lifecycle hooks:
```ts
watch(
  () => props.slug,                 // like [slug]
  async (slug, _old, onCleanup) => {
    const controller = new AbortController();
    onCleanup(() => controller.abort());
    data.value = await api.getBySlug(slug, controller.signal);
  },
  { immediate: true },              // also run on mount, like useEffect does
);
```
Key difference: Vue tracks reactive dependencies automatically, while React makes you list them. The `react-hooks/exhaustive-deps` lint rule (enabled in `eslint.config.js`) warns when you forget one.

---

## 12. App Bootstrap: `createRoot` & `StrictMode`
### React
```tsx
// main.tsx
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
```
`StrictMode` renders nothing visible. In development only, it runs extra checks (double rendering, double effects) to surface impure components and missing cleanups.

### Vue 3 Equivalent: `createApp().mount()`
```ts
createApp(App).use(router).mount('#root');
```
Vue has no StrictMode. Note that in React the router is a component (`<RouterProvider>`, see 14), not a plugin installed with `.use()`.

---

## 13. Custom Hooks
### React: functions named `useSomething` that call other hooks
A custom hook packages stateful logic for reuse. It's just a function, but because it calls hooks it must follow the **Rules of Hooks**:
1. Only call hooks at the top level of a component or another hook (never inside `if`, loops or callbacks), because React identifies each hook by its call order.
2. Only call hooks from React components or other hooks.
```ts
// hooks/useApi.ts
export function useApi<T>(fetcher: (signal: AbortSignal) => Promise<T>) {
  const [state, setState] = useState({ data: null, loading: true, error: null });
  useEffect(() => { /* fetch, setState, abort on cleanup */ }, [fetcher]);
  return state;
}

// pages/HomePage.tsx
const { data: accommodations, loading, error } = useApi(api.listAccommodations);
```
Each component that calls `useApi` gets its **own** independent state. Hooks share logic, not state.

### Vue 3 Equivalent: Composables
Same idea and naming convention (`useFetch`, `useMouse`). Vue composables have no call-order rule, since `ref`s are real objects rather than slots looked up by position, and you can call them conditionally inside `setup`.
```ts
export function useApi<T>(fetcher: (signal: AbortSignal) => Promise<T>) {
  const data = ref<T | null>(null);
  const loading = ref(true);
  const error = ref<string | null>(null);
  const controller = new AbortController();
  fetcher(controller.signal).then(d => (data.value = d)).catch(...).finally(() => (loading.value = false));
  onUnmounted(() => controller.abort());
  return { data, loading, error };
}
```

---

## 14. Routing: React Router vs Vue Router
### React: React Router 7
Routes are declared as a config array, much like Vue Router.
```tsx
// router.tsx
export const routes: RouteObject[] = [
  {
    path: '/',
    element: <SiteLayout />,               // shared navbar + footer
    children: [
      { index: true, element: <HomePage /> },   // "/" itself
      { path: '*', element: <NotFoundPage /> }, // catch-all 404
    ],
  },
];
export const router = createBrowserRouter(routes);

// App.tsx
<RouterProvider router={router} />

// layouts/SiteLayout.tsx: child routes render where <Outlet /> is
<Navbar />
<Outlet />
<Footer />

// Navigation without a page reload
<Link to="/">Return to the Lobby</Link>
```

### Vue 3 Equivalent: Vue Router 4
| React Router | Vue Router |
|---|---|
| `createBrowserRouter(routes)` | `createRouter({ history: createWebHistory(), routes })` |
| `<RouterProvider router={router} />` | `app.use(router)` + `<router-view />` in `App.vue` |
| `element: <HomePage />` | `component: HomePage` |
| `{ index: true, ... }` | child with `path: ''` |
| `<Outlet />` | nested `<router-view />` |
| `<Link to="/">` | `<router-link to="/">` |
| `path: '*'` | `path: '/:pathMatch(.*)*'` |
| `useParams()` / `useNavigate()` | `useRoute().params` / `useRouter().push()` |
| `createMemoryRouter` (tests) | `createMemoryHistory()` |

Note that `element` takes an already-created JSX element (`<HomePage />`), whereas Vue's `component` takes the component itself (`HomePage`).

---

## 15. Component Tests
### React: Vitest + React Testing Library
Tests render components into a simulated DOM (jsdom) and query it the way a user would, by visible text and accessible roles, rather than by component internals.
```tsx
// components/MerchShop.test.tsx
mockFetch({ '/merchandise': products });
render(<MerchShop />);

expect(await screen.findByText('Gotham Skyline Silk Tie')).toBeInTheDocument(); // findBy* waits for async
await userEvent.click(screen.getByRole('button', { name: 'apparel' }));
expect(screen.queryByText('Kane County Reserve Wine')).not.toBeInTheDocument();  // queryBy* returns null instead of throwing
```
Run with `npm test` (once) or `npm run test:watch`.

### Vue 3 Equivalent: Vitest + Vue Test Utils (or Vue Testing Library)
```ts
const wrapper = mount(MerchShop);
await flushPromises();
expect(wrapper.text()).toContain('Gotham Skyline Silk Tie');
await wrapper.find('button.filter-btn').trigger('click');
```
`@testing-library/vue` also exists with the same `render`/`screen` API as above, so the React style carries over directly.
