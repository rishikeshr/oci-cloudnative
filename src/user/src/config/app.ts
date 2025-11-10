
export interface IAppEnv {
  POSTGRES_HOST: string;
  POSTGRES_PORT: string;
  POSTGRES_DB: string;
  POSTGRES_USER: string;
  POSTGRES_PASSWORD: string;
  NODE_ENV?: string;
}

export type IDBConfig = Pick<IAppEnv, 'POSTGRES_HOST' | 'POSTGRES_PORT' | 'POSTGRES_DB' | 'POSTGRES_USER' | 'POSTGRES_PASSWORD'>;

export class AppConfig {
  private static COMMON: AppConfig = new AppConfig();
  public static common(c?: AppConfig): AppConfig {
    this.COMMON = c || this.COMMON;
    return this.COMMON;
  }

  constructor(private ENV: IAppEnv = process.env as any) {

  }

  public env(): IAppEnv {
    return {...this.ENV};
  }

  public prod(): boolean {
    const { NODE_ENV } = this.ENV;
    return /^prod/i.test(NODE_ENV || '');
  }

  /**
   * get db configuration
   */
  public dbConfig(): IDBConfig {
    const { POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD } = this.ENV;
    return {
      POSTGRES_HOST,
      POSTGRES_PORT,
      POSTGRES_DB,
      POSTGRES_USER,
      POSTGRES_PASSWORD,
    };
  }

  /**
   * check if should use mock (memory) db
   */
  public mockDb(): boolean {
    // tslint:disable-next-line:triple-equals
    return this.dbConfig().POSTGRES_HOST == 'mock';
  }

}
