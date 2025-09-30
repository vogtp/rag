import { HttpHeaders } from '@angular/common/http';

export const collectionURL = '/vecdb/';
export const summaryURL = '/summary/';
export const userSettingsURL = '/user/';

export const httpHeaders: HttpHeaders = new HttpHeaders({
  //Authorization: 'Bearer JWT-token'
  Accept: 'application/json',
});
