import { useEffect, useState } from 'react';

interface ApiState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

// Runs `fetcher` when the component mounts and tracks its loading/error state.
// The request is aborted if the component unmounts before it finishes.
//
// `fetcher` must keep the same identity between renders (e.g. a function from
// `api` in api/client.ts); a new inline arrow each render would refetch forever.
export function useApi<T>(fetcher: (signal: AbortSignal) => Promise<T>): ApiState<T> {
  const [state, setState] = useState<ApiState<T>>({ data: null, loading: true, error: null });

  useEffect(() => {
    const controller = new AbortController();

    fetcher(controller.signal)
      .then((data) => setState({ data, loading: false, error: null }))
      .catch((err: unknown) => {
        if (controller.signal.aborted) return; // unmounted: nobody to tell
        console.error(err);
        setState({
          data: null,
          loading: false,
          error: err instanceof Error ? err.message : 'Something went wrong',
        });
      });

    return () => controller.abort();
  }, [fetcher]);

  return state;
}
