import type { Client } from '@urql/svelte';
import { createClient, fetchExchange } from '@urql/svelte';
import { cacheExchange } from '@urql/exchange-graphcache';
import { persistedExchange } from '@urql/exchange-persisted';

import schema from '$lib/generated/graphql.schema.urql.json';

export function initializeGraphQLClient(): Client {
  return createClient({
    url: 'https://api.ficsit.app/v2/query',
    exchanges: [
      cacheExchange({
        schema,
        keys: {
          GetMods: () => null,
          GetSMLVersions: () => null,
          LatestVersions: () => null,
          UserMod: () => null,
          GetGuides: () => null,
          OAuthOptions: () => null,
          UserRoles: () => null,
          Compatibility: () => null,
          CompatibilityInfo: () => null,
          VersionDependency: () => null,
          Mod: (data) => data.mod_reference,
        },
        resolvers: {
          Query: {
            getModByReference: (_parent, args) => {
              return { __typename: 'Mod', mod_reference: args.modReference };
            },
          },
        },
      }),
      persistedExchange({
        preferGetForPersistedQueries: true,
      }),
      fetchExchange,
    ],
  });
}
