import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from './environment';
import { RunRequest } from './models/Run';

@Injectable({
  providedIn: 'root'
})
export class RunService {

  constructor(private http: HttpClient) { }

  run(req: RunRequest): Observable<any> {
    return this.http.post(`${environment.apiUrl}/run`, req);
  }
  runOutput(): Observable<string> {
    return this.http.get(`${environment.apiUrl}/run/output`, {responseType: 'text'});
  }
}
