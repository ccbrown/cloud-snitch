import { useCallback } from 'react';

import { MapLocation, MapRect } from './map';
import { useSearchParamState } from '@/hooks';

export type Selection =
    | {
          type: 'aws-region';
          id: string;
      }
    | {
          type: 'network';
          cidr: string;
      }
    | {
          type: 'principal';
          id: string;
      }
    | {
          type: 'cluster';
          rect: MapRect;
          location: MapLocation;
      };

export const isEqualSelection = (a: Selection | null, b: Selection | null): boolean => {
    if (!a || !b) {
        return a === b;
    }
    switch (a.type) {
        case 'aws-region':
            return a.type === b.type && a.id === b.id;
        case 'network':
            return a.type === b.type && a.cidr === b.cidr;
        case 'principal':
            return a.type === b.type && a.id === b.id;
        case 'cluster':
            return a.type === b.type && a.rect.equals(b.rect);
        default:
            return false;
    }
};

export const parseSelection = (selection: string | null): Selection | null => {
    if (!selection) {
        return null;
    }

    const colon = selection.indexOf(':');
    if (colon === -1) {
        return null;
    }

    const parts = [selection.slice(0, colon), selection.slice(colon + 1)];

    switch (parts[0]) {
        case 'aws-region':
            return { type: 'aws-region', id: parts[1] };
        case 'network':
            return { type: 'network', cidr: parts[1] };
        case 'principal':
            return { type: 'principal', id: parts[1] };
        case 'cluster':
            const rectParts = parts[1].split(',');
            if (rectParts.length !== 6) {
                return null;
            }
            return {
                type: 'cluster',
                rect: new MapRect({
                    minMercatorX: parseFloat(rectParts[0]),
                    minMercatorY: parseFloat(rectParts[1]),
                    maxMercatorX: parseFloat(rectParts[2]),
                    maxMercatorY: parseFloat(rectParts[3]),
                }),
                location: MapLocation.fromLatitudeAndLongitude(parseFloat(rectParts[4]), parseFloat(rectParts[5])),
            };
        default:
            return null;
    }
};

export const stringifySelection = (selection: Selection | null): string | null => {
    if (!selection) {
        return null;
    }

    switch (selection.type) {
        case 'aws-region':
            return `aws-region:${selection.id}`;
        case 'network':
            return `network:${selection.cidr}`;
        case 'principal':
            return `principal:${selection.id}`;
        case 'cluster':
            return `cluster:${selection.rect.minMercatorX},${selection.rect.minMercatorY},${selection.rect.maxMercatorX},${selection.rect.maxMercatorY},${selection.location.latitude},${selection.location.longitude}`;
        default:
            return null;
    }
};

export const useSelection = (): [Selection | null, (selection: Selection | null) => void] => {
    const [selectionString, setSelectionString] = useSearchParamState('selection');
    const selection = parseSelection(selectionString);
    const setSelection = useCallback(
        (selection: Selection | null) => {
            setSelectionString(stringifySelection(selection) || '');
        },
        [setSelectionString],
    );
    return [selection, setSelection];
};
